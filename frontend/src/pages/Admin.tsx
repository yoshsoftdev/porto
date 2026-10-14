import { useEffect, useState, type FormEvent } from "react";
import { Link } from "react-router-dom";
import { asset, getJSON, sendJSON, uploadImage } from "../api";
import type { LLMConfig, Profile, Project, Section, SkillGroup } from "../types";

const TOKEN_KEY = "porto_admin_token";

export default function Admin() {
  const [token, setToken] = useState(sessionStorage.getItem(TOKEN_KEY) || "");
  const [password, setPassword] = useState("");
  const [tab, setTab] = useState<"profile" | "projects" | "skills" | "sections" | "llm">("profile");
  const [error, setError] = useState("");
  const [ok, setOk] = useState("");

  async function login(e: FormEvent) {
    e.preventDefault();
    setError("");
    try {
      const res = await sendJSON<{ token: string }>("/api/admin/login", "POST", { password });
      sessionStorage.setItem(TOKEN_KEY, res.token);
      setToken(res.token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "login failed");
    }
  }

  if (!token) {
    return (
      <form className="login" onSubmit={login}>
        <div className="kicker">Restricted</div>
        <h1>Admin</h1>
        <p style={{ color: "var(--muted)" }}>Password-protected CRUD for profile, projects, images, and LLM keys.</p>
        <label>Password</label>
        <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        {error && <div className="flash err">{error}</div>}
        <div className="actions">
          <button className="primary" type="submit">
            Enter
          </button>
          <Link className="ghost" to="/" style={{ textDecoration: "none", display: "inline-block" }}>
            Back to site
          </Link>
        </div>
      </form>
    );
  }

  return (
    <div className="admin">
      <div className="topbar-inner" style={{ padding: "0 0 12px" }}>
        <div>
          <div className="kicker">Admin</div>
          <h1 style={{ fontSize: 42, margin: "6px 0" }}>Edit portfolio</h1>
        </div>
        <div className="actions">
          <Link className="ghost" to="/" style={{ textDecoration: "none" }}>
            View site
          </Link>
          <button
            className="danger"
            onClick={() => {
              sessionStorage.removeItem(TOKEN_KEY);
              setToken("");
            }}
          >
            Logout
          </button>
        </div>
      </div>
      {ok && <div className="flash ok">{ok}</div>}
      {error && <div className="flash err">{error}</div>}
      <div className="tabs">
        {(["profile", "projects", "skills", "sections", "llm"] as const).map((t) => (
          <button key={t} className="ghost" aria-selected={tab === t} onClick={() => setTab(t)}>
            {t}
          </button>
        ))}
      </div>
      {tab === "profile" && <ProfilePane token={token} onError={setError} onOk={setOk} />}
      {tab === "projects" && <ProjectsPane token={token} onError={setError} onOk={setOk} />}
      {tab === "skills" && <SkillsPane token={token} onError={setError} onOk={setOk} />}
      {tab === "sections" && <SectionsPane token={token} onError={setError} onOk={setOk} />}
      {tab === "llm" && <LLMPane token={token} onError={setError} onOk={setOk} />}
    </div>
  );
}

function ProfilePane({
  token,
  onError,
  onOk,
}: {
  token: string;
  onError: (s: string) => void;
  onOk: (s: string) => void;
}) {
  const [p, setP] = useState<Profile | null>(null);

  useEffect(() => {
    getJSON<Profile>("/api/admin/profile", token)
      .then(setP)
      .catch((e) => onError(e.message));
  }, [token, onError]);

  if (!p) return null;

  async function save() {
    onError("");
    try {
      const saved = await sendJSON<Profile>("/api/admin/profile", "PUT", p, token);
      setP((cur) => ({ ...saved, photo_url: cur?.photo_url, skills: cur?.skills || [], sections: cur?.sections || [] }));
      onOk("Profile saved");
    } catch (e) {
      onError(e instanceof Error ? e.message : "save failed");
    }
  }

  async function onPhoto(file: File) {
    if (!p) return;
    try {
      const img = await uploadImage(file, token);
      const next: Profile = { ...p, photo_id: img.id, photo_url: img.url };
      setP(next);
      await sendJSON("/api/admin/profile", "PUT", next, token);
      onOk("Photo updated");
    } catch (e) {
      onError(e instanceof Error ? e.message : "upload failed");
    }
  }

  return (
    <div className="card">
      <div className="row">
        <div>
          <label>Name</label>
          <input value={p.name} onChange={(e) => setP({ ...p, name: e.target.value })} />
        </div>
        <div>
          <label>Title</label>
          <input value={p.title} onChange={(e) => setP({ ...p, title: e.target.value })} />
        </div>
      </div>
      <label>Target role</label>
      <input value={p.target_role} onChange={(e) => setP({ ...p, target_role: e.target.value })} />
      <label>Summary</label>
      <textarea style={{ minHeight: 160 }} value={p.summary} onChange={(e) => setP({ ...p, summary: e.target.value })} />
      <div className="row">
        <div>
          <label>Email</label>
          <input value={p.email} onChange={(e) => setP({ ...p, email: e.target.value })} />
        </div>
        <div>
          <label>GitHub</label>
          <input value={p.github} onChange={(e) => setP({ ...p, github: e.target.value })} />
        </div>
      </div>
      <label>LinkedIn</label>
      <input value={p.linkedin} onChange={(e) => setP({ ...p, linkedin: e.target.value })} />
      <label>Strengths (one per line)</label>
      <textarea
        value={(p.strengths || []).join("\n")}
        onChange={(e) => setP({ ...p, strengths: e.target.value.split("\n").map((s) => s.trim()).filter(Boolean) })}
      />
      <label>Profile photo</label>
      {p.photo_url && <img src={asset(p.photo_url)} alt="" style={{ width: 160, marginBottom: 8, border: "1px solid var(--line)" }} />}
      <input type="file" accept="image/*" onChange={(e) => e.target.files?.[0] && void onPhoto(e.target.files[0])} />
      <div className="actions">
        <button className="primary" onClick={() => void save()}>
          Save profile
        </button>
      </div>
    </div>
  );
}

function emptyProject(): Project {
  return {
    id: 0,
    slug: "",
    name: "",
    subtitle: "",
    period: "",
    client: "",
    status: "",
    role: "",
    team: "",
    lead: "",
    extra: [],
    summary: "",
    tech: [],
    content: "",
    sort_order: 0,
    hidden: false,
    images: [],
    image_urls: [],
  };
}

function ProjectsPane({
  token,
  onError,
  onOk,
}: {
  token: string;
  onError: (s: string) => void;
  onOk: (s: string) => void;
}) {
  const [list, setList] = useState<Project[]>([]);
  const [cur, setCur] = useState<Project>(emptyProject());

  async function reload() {
    const items = await getJSON<Project[]>("/api/admin/projects", token);
    setList(items);
  }

  useEffect(() => {
    reload().catch((e) => onError(e.message));
  }, [token]);

  async function save() {
    onError("");
    try {
      if (cur.id) {
        const saved = await sendJSON<Project>(`/api/admin/projects/${cur.id}`, "PUT", cur, token);
        setCur(saved);
      } else {
        const created = await sendJSON<Project>("/api/admin/projects", "POST", cur, token);
        setCur(created);
      }
      await reload();
      onOk("Project saved");
    } catch (e) {
      onError(e instanceof Error ? e.message : "save failed");
    }
  }

  async function remove() {
    if (!cur.id || !confirm("Delete this project?")) return;
    await sendJSON(`/api/admin/projects/${cur.id}`, "DELETE", undefined, token);
    setCur(emptyProject());
    await reload();
    onOk("Project deleted");
  }

  async function addShot(file: File) {
    if (!cur.id) {
      onError("Save the project first, then add screenshots");
      return;
    }
    const img = await uploadImage(file, token);
    await sendJSON(`/api/admin/projects/${cur.id}/images`, "POST", { image_id: img.id }, token);
    await reload();
    const fresh = await getJSON<Project[]>("/api/admin/projects", token);
    setCur(fresh.find((p) => p.id === cur.id) || cur);
    onOk("Screenshot added");
  }

  async function detach(relId: number) {
    await sendJSON(`/api/admin/projects/${cur.id}/images/${relId}`, "DELETE", undefined, token);
    const fresh = await getJSON<Project[]>("/api/admin/projects", token);
    setList(fresh);
    setCur(fresh.find((p) => p.id === cur.id) || cur);
  }

  return (
    <div className="row" style={{ alignItems: "start" }}>
      <div className="card">
        <div className="actions" style={{ marginTop: 0 }}>
          <button className="primary" onClick={() => setCur(emptyProject())}>
            New project
          </button>
        </div>
        <div className="proj-list">
          {list.map((p) => (
            <button key={p.id} className="ghost" onClick={() => setCur(p)}>
              {p.sort_order}. {p.name}
            </button>
          ))}
        </div>
      </div>
      <div className="card">
        <div className="row">
          <div>
            <label>Name</label>
            <input value={cur.name} onChange={(e) => setCur({ ...cur, name: e.target.value })} />
          </div>
          <div>
            <label>Subtitle</label>
            <input value={cur.subtitle} onChange={(e) => setCur({ ...cur, subtitle: e.target.value })} />
          </div>
        </div>
        <div className="row">
          <div>
            <label>Slug</label>
            <input value={cur.slug} onChange={(e) => setCur({ ...cur, slug: e.target.value })} />
          </div>
          <div>
            <label>Period</label>
            <input value={cur.period} onChange={(e) => setCur({ ...cur, period: e.target.value })} />
          </div>
        </div>
        <div className="row">
          <div>
            <label>Client</label>
            <input value={cur.client} onChange={(e) => setCur({ ...cur, client: e.target.value })} />
          </div>
          <div>
            <label>Status</label>
            <input value={cur.status} onChange={(e) => setCur({ ...cur, status: e.target.value })} />
          </div>
        </div>
        <div className="row">
          <div>
            <label>Role</label>
            <input value={cur.role} onChange={(e) => setCur({ ...cur, role: e.target.value })} />
          </div>
          <div>
            <label>Sort</label>
            <input type="number" value={cur.sort_order} onChange={(e) => setCur({ ...cur, sort_order: Number(e.target.value) })} />
          </div>
        </div>
        <label>Summary</label>
        <textarea value={cur.summary} onChange={(e) => setCur({ ...cur, summary: e.target.value })} />
        <label>Tech (comma separated)</label>
        <input value={(cur.tech || []).join(", ")} onChange={(e) => setCur({ ...cur, tech: e.target.value.split(",").map((s) => s.trim()).filter(Boolean) })} />
        <label>Content (markdown)</label>
        <textarea style={{ minHeight: 220 }} value={cur.content} onChange={(e) => setCur({ ...cur, content: e.target.value })} />
        <label>
          <input type="checkbox" checked={cur.hidden} onChange={(e) => setCur({ ...cur, hidden: e.target.checked })} /> Hidden
        </label>
        <label>Screenshots</label>
        <div className="thumb-row">
          {(cur.images || []).map((img) => (
            <div className="thumb" key={img.id}>
              <img src={asset(`/api/images/${img.image_id}`)} alt="" />
              <button className="danger" onClick={() => void detach(img.id)}>
                Remove
              </button>
            </div>
          ))}
        </div>
        <input type="file" accept="image/*" onChange={(e) => e.target.files?.[0] && void addShot(e.target.files[0])} />
        <div className="actions">
          <button className="primary" onClick={() => void save()}>
            Save project
          </button>
          {cur.id > 0 && (
            <button className="danger" onClick={() => void remove()}>
              Delete
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

function SkillsPane({
  token,
  onError,
  onOk,
}: {
  token: string;
  onError: (s: string) => void;
  onOk: (s: string) => void;
}) {
  const [items, setItems] = useState<SkillGroup[]>([]);
  useEffect(() => {
    getJSON<SkillGroup[]>("/api/admin/skills", token)
      .then(setItems)
      .catch((e) => onError(e.message));
  }, [token, onError]);

  return (
    <div>
      {items.map((g, i) => (
        <div className="card" key={i}>
          <label>Category</label>
          <input
            value={g.category}
            onChange={(e) => {
              const next = [...items];
              next[i] = { ...g, category: e.target.value };
              setItems(next);
            }}
          />
          <label>Items (one per line)</label>
          <textarea
            value={g.items.join("\n")}
            onChange={(e) => {
              const next = [...items];
              next[i] = { ...g, items: e.target.value.split("\n").map((s) => s.trim()).filter(Boolean) };
              setItems(next);
            }}
          />
          <button className="danger" onClick={() => setItems(items.filter((_, j) => j !== i))}>
            Remove group
          </button>
        </div>
      ))}
      <div className="actions">
        <button className="ghost" onClick={() => setItems([...items, { category: "", items: [], sort_order: items.length + 1 }])}>
          Add group
        </button>
        <button
          className="primary"
          onClick={() =>
            sendJSON("/api/admin/skills", "PUT", items, token)
              .then(() => onOk("Skills saved"))
              .catch((e) => onError(e.message))
          }
        >
          Save skills
        </button>
      </div>
    </div>
  );
}

function SectionsPane({
  token,
  onError,
  onOk,
}: {
  token: string;
  onError: (s: string) => void;
  onOk: (s: string) => void;
}) {
  const [items, setItems] = useState<Section[]>([]);
  useEffect(() => {
    getJSON<Section[]>("/api/admin/sections", token)
      .then(setItems)
      .catch((e) => onError(e.message));
  }, [token, onError]);

  return (
    <div>
      {items.map((s, i) => (
        <div className="card" key={i}>
          <label>Title</label>
          <input
            value={s.title}
            onChange={(e) => {
              const next = [...items];
              next[i] = { ...s, title: e.target.value };
              setItems(next);
            }}
          />
          <label>Content</label>
          <textarea
            style={{ minHeight: 160 }}
            value={s.content}
            onChange={(e) => {
              const next = [...items];
              next[i] = { ...s, content: e.target.value };
              setItems(next);
            }}
          />
          <label>
            <input
              type="checkbox"
              checked={s.hidden}
              onChange={(e) => {
                const next = [...items];
                next[i] = { ...s, hidden: e.target.checked };
                setItems(next);
              }}
            />{" "}
            Hidden
          </label>
        </div>
      ))}
      <button
        className="primary"
        onClick={() =>
          sendJSON("/api/admin/sections", "PUT", items, token)
            .then(() => onOk("Sections saved"))
            .catch((e) => onError(e.message))
        }
      >
        Save sections
      </button>
    </div>
  );
}

function LLMPane({
  token,
  onError,
  onOk,
}: {
  token: string;
  onError: (s: string) => void;
  onOk: (s: string) => void;
}) {
  const [cfg, setCfg] = useState<LLMConfig | null>(null);
  const [keys, setKeys] = useState<Record<string, string>>({});
  useEffect(() => {
    getJSON<LLMConfig>("/api/admin/llm", token)
      .then(setCfg)
      .catch((e) => onError(e.message));
  }, [token, onError]);
  if (!cfg) return null;
  const labels: Record<string, string> = { gemini: "Gemini", gpt: "GPT", grok: "Grok", claude: "Claude" };

  return (
    <div className="card">
      <p style={{ color: "var(--muted)" }}>
        Store provider keys here. The public chatbot can switch among providers that have a key. Answers stay grounded in porto.md.
      </p>
      <label>Default provider</label>
      <select value={cfg.active} onChange={(e) => setCfg({ ...cfg, active: e.target.value })}>
        {cfg.providers.map((p) => (
          <option key={p.id} value={p.id}>
            {labels[p.id]}
          </option>
        ))}
      </select>
      {cfg.providers.map((p, i) => (
        <div key={p.id} className="card" style={{ marginTop: 12 }}>
          <strong>{labels[p.id]}</strong>
          <label>Model</label>
          <input
            value={p.model}
            onChange={(e) => {
              const providers = [...cfg.providers];
              providers[i] = { ...p, model: e.target.value };
              setCfg({ ...cfg, providers });
            }}
          />
          <label>API key {p.has_key ? "(already saved — leave blank to keep)" : ""}</label>
          <input type="password" value={keys[p.id] || ""} onChange={(e) => setKeys({ ...keys, [p.id]: e.target.value })} />
        </div>
      ))}
      <div className="actions">
        <button
          className="primary"
          onClick={() =>
            sendJSON(
              "/api/admin/llm",
              "PUT",
              {
                active: cfg.active,
                providers: cfg.providers.map((p) => ({ id: p.id, model: p.model, key: keys[p.id] || "" })),
              },
              token,
            )
              .then((next) => {
                setCfg(next as LLMConfig);
                setKeys({});
                onOk("LLM settings saved");
              })
              .catch((e) => onError(e.message))
          }
        >
          Save LLM settings
        </button>
      </div>
    </div>
  );
}
