-- name: CreateMembers :copyfrom
INSERT INTO members (group_id, name) VALUES (@group_id, @name);
