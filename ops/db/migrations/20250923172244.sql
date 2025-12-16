-- Create policy "group_policy"
CREATE POLICY "group_policy" ON "members" AS PERMISSIVE FOR ALL TO PUBLIC USING (group_id = (current_setting('app.group_id'::text))::integer);
-- Modify "members" table
ALTER TABLE "members" ENABLE ROW LEVEL SECURITY;
