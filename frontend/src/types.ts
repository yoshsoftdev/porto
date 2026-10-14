export type MetaItem = { key: string; value: string };

export type Profile = {
  id: number;
  name: string;
  title: string;
  summary: string;
  target_role: string;
  email: string;
  github: string;
  linkedin: string;
  strengths: string[];
  photo_id?: number | null;
  photo_url?: string;
  skills: SkillGroup[];
  sections: Section[];
};

export type SkillGroup = {
  id?: number;
  category: string;
  items: string[];
  sort_order: number;
};

export type Section = {
  id?: number;
  slug: string;
  title: string;
  content: string;
  sort_order: number;
  hidden: boolean;
};

export type ProjectImage = {
  id: number;
  project_id: number;
  image_id: number;
  caption: string;
  sort_order: number;
};

export type Project = {
  id: number;
  slug: string;
  name: string;
  subtitle: string;
  period: string;
  client: string;
  status: string;
  role: string;
  team: string;
  lead: string;
  extra: MetaItem[] | null;
  summary: string;
  tech: string[] | null;
  content: string;
  sort_order: number;
  hidden: boolean;
  images: ProjectImage[] | null;
  image_urls: string[];
};

export type ChatProvider = {
  id: string;
  label: string;
  enabled: boolean;
};

export type LLMConfig = {
  active: string;
  providers: { id: string; model: string; has_key: boolean; enabled: boolean }[];
};
