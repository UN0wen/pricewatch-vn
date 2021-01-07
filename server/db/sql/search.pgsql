CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE OR REPLACE FUNCTION public.immutable_unaccent (regdictionary, text)
  RETURNS text
  LANGUAGE c
  IMMUTABLE STRICT
  AS '$libdir/unaccent'
,
  'unaccent_dict';

CREATE OR REPLACE FUNCTION public.f_lower_unaccent (text)
  RETURNS text
  LANGUAGE sql
  IMMUTABLE STRICT
  AS $func$
  SELECT
    lower(public.immutable_unaccent (regdictionary 'public.unaccent', $1))
$func$;

-- CREATE INDEX items_unaccent_name_idx ON items (public.f_unaccent (name));
-- CREATE INDEX items_unaccent_name_trgm_idx ON items USING gin (public.f_unaccent (name) gin_trgm_ops);
CREATE INDEX items_lower_unaccent_name_trgm_idx2 ON items USING gin (f_lower_unaccent (name) gin_trgm_ops);



