package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	url string
}

func NewRedisRepository(redisURL string) (*RedisRepository, error) {
	if _, err := url.Parse(redisURL); err != nil {
		return nil, err
	}

	return &RedisRepository{url: redisURL}, nil
}

func (r *RedisRepository) client() (*redis.Client, error) {
	clientOptions, err := redis.ParseURL(r.url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(clientOptions)
	return client, nil
}

func (r *RedisRepository) Ping() (string, error) {
	client, err := r.client()
	if err != nil {
		return "", err
	}
	defer client.Close()

	return client.Ping(context.Background()).Result()
}

func (r *RedisRepository) CreateMessage(title, body string) (string, error) {
	messageID := uuid.NewString()

	message := struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		ID:    messageID,
		Title: title,
		Body:  body,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	client, err := r.client()
	if err != nil {
		return "", err
	}
	defer client.Close()

	if err := client.RPush(context.Background(), "messageboard", data).Err(); err != nil {
		return "", err
	}

	return messageID, nil
}

func (r *RedisRepository) ListMessages() (string, error) {
	return r.listAsJSON("messageboard")
}

func (r *RedisRepository) CreateReply(messageID, reply string) error {
	replyJSON := struct {
		Reply string `json:"reply"`
	}{Reply: reply}

	data, err := json.Marshal(replyJSON)
	if err != nil {
		return err
	}

	client, err := r.client()
	if err != nil {
		return err
	}
	defer client.Close()

	return client.RPush(context.Background(), "messageboard:"+messageID, data).Err()
}

func (r *RedisRepository) ListReplies(messageID string) (string, error) {
	return r.listAsJSON("messageboard:" + messageID)
}

func (r *RedisRepository) listAsJSON(key string) (string, error) {
	client, err := r.client()
	if err != nil {
		return "", err
	}
	defer client.Close()

	entries, err := client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return "", err
	}

	result := make([]json.RawMessage, 0, len(entries))
	for _, entry := range entries {
		if !json.Valid([]byte(entry)) {
			return "", fmt.Errorf("invalid JSON stored in Redis key %q: %v", key, entry)
		}
		result = append(result, json.RawMessage(entry))
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

var _ Repository = (*RedisRepository)(nil)
