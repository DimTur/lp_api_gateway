package redis

import (
	"context"
	"fmt"
	"time"
)

func (r *RedisClient) SaveLgUser(ctx context.Context, userID string, groupIDs []string) error {
	const op = "storage.redis.SaveLgUser"

	key := fmt.Sprintf("user_groups:%s", userID)
	members := make([]interface{}, len(groupIDs))
	for i, groupID := range groupIDs {
		members[i] = groupID
	}

	if err := r.client.SAdd(ctx, key, members...).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := r.client.Expire(ctx, key, 1*time.Minute).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisClient) SaveLgShareWithChannel(ctx context.Context, channelID int64, groupIDs []string) error {
	const op = "storage.redis.SaveLgShareWithChannel"

	key := fmt.Sprintf("channel_shared_groups:%d", channelID)
	members := make([]interface{}, len(groupIDs))
	for i, groupID := range groupIDs {
		members[i] = groupID
	}

	if err := r.client.SAdd(ctx, key, members...).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := r.client.Expire(ctx, key, 1*time.Minute).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisClient) CheckGroupsIntersection(ctx context.Context, userID string, channelID int64) (bool, error) {
	const op = "storage.redis.CheckGroupsIntersection"

	userKey := fmt.Sprintf("user_groups:%s", userID)
	channelKey := fmt.Sprintf("channel_shared_groups:%d", channelID)

	intersection, err := r.client.SInter(ctx, userKey, channelKey).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.client.Del(ctx, userKey, channelKey).Err(); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return len(intersection) > 0, nil
}
