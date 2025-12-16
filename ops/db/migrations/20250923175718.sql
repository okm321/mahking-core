-- Drop "on_update_time_members" trigger
DROP TRIGGER "on_update_time_members" ON "groups";
-- Create trigger "on_update_time_members"
CREATE TRIGGER "on_update_time_members" BEFORE UPDATE ON "members" FOR EACH ROW EXECUTE FUNCTION "update_updated_at_column"();
