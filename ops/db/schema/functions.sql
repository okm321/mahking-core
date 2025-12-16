CREATE OR REPLACE FUNCTION update_updated_at_column()returns TRIGGER AS
$$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$
language plpgsql;
