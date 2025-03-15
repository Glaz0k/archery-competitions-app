ALTER TABLE "competitions" DROP CONSTRAINT IF EXISTS "competitions_cup_id_fkey";
ALTER TABLE "individual_groups" DROP CONSTRAINT IF EXISTS "individual_groups_competition_id_fkey";

DROP TABLE IF EXISTS "competitors_detail";
DROP TABLE IF EXISTS "competitors";
DROP TABLE IF EXISTS "shots";
DROP TABLE IF EXISTS "ranges";
DROP TABLE IF EXISTS "range_groups";
DROP TABLE IF EXISTS "shoot_outs";
DROP TABLE IF EXISTS "sparring_places";
DROP TABLE IF EXISTS "sparrings";
DROP TABLE IF EXISTS "finals";
DROP TABLE IF EXISTS "semifinals";
DROP TABLE IF EXISTS "quarterfinals";
DROP TABLE IF EXISTS "qualification_rounds";
DROP TABLE IF EXISTS "qualification_sections";
DROP TABLE IF EXISTS "qualifications";
DROP TABLE IF EXISTS "competitor_group_details";
DROP TABLE IF EXISTS "competitor_competition_details";
DROP TABLE IF EXISTS "individual_groups";
DROP TABLE IF EXISTS "competitions";
DROP TABLE IF EXISTS "cups";

DROP TYPE IF EXISTS "group_state";
DROP TYPE IF EXISTS "gender";
DROP TYPE IF EXISTS "bow_class";
DROP TYPE IF EXISTS "sports_rank";
DROP TYPE IF EXISTS "sparring_state";
DROP TYPE IF EXISTS "competition_stage";