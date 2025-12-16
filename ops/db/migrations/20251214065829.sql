-- Create index "idx_members_group_id" to table: "members"
CREATE INDEX "idx_members_group_id" ON "members" ("group_id");
-- Create index "idx_rules_group_id" to table: "rules"
CREATE INDEX "idx_rules_group_id" ON "rules" ("group_id");
