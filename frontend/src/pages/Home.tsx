import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { asset, getJSON } from "../api";
import Chatbot from "../components/Chatbot";
import Markdown from "../components/Markdown";
import type { Profile, Project } from "../types";

export default function Home() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [projects, setProjects] = useState<Project[]>([]);
  const [open, setOpen] = useState<string | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    Promise.all([getJSON<Profile>("/api/profile"), getJSON<Project[]>("/api/projects")])
      .then(([p, list]) => {
        setProfile(p);
        setProjects(list);
        setOpen(list[0]?.slug ?? null);
      })
      .catch((e) => setError(e instanceof Error ? e.message : "failed to load"));
  }, []);

  return (
    <>
      <header className="topbar">
        <div className="wrap topbar-inner">
          <a className="brand" href="#/">
            <b>Yo Shandy</b>
            <span>Backend & Integration</span>
          </a>
          <nav className="nav">
            <a
              href="#/"
              onClick={(e) => {
                e.preventDefault();
                document.getElementById("projects")?.scrollIntoView({ behavior: "smooth" });
              }}
            >
              Projects
            </a>
            <a
              href="#/"
              onClick={(e) => {
                e.preventDefault();
                document.getElementById("approach")?.scrollIntoView({ behavior: "smooth" });
              }}
            >
              Approach
            </a>
            <a
              href="#/"
              onClick={(e) => {
                e.preventDefault();
                document.getElementById("contact")?.scrollIntoView({ behavior: "smooth" });
              }}
            >
              Contact
            </a>
            <Link className="admin-link" to="/admin">
              Admin
            </Link>
          </nav>
        </div>
      </header>

      <main className="wrap">
        {error && <div className="flash err">{error}</div>}
        {profile && (
          <>
            <section className="hero">
              <div>
                <div className="kicker">{profile.target_role || profile.title}</div>
                <h1>{profile.name}</h1>
                <p className="lede">{profile.summary}</p>
                <div className="meta-row">
                  {profile.email && <span className="chip">{profile.email}</span>}
                  {profile.github && (
                    <a className="chip" href={profile.github} target="_blank" rel="noreferrer">
                      GitHub
                    </a>
                  )}
                  {profile.linkedin && (
                    <a className="chip" href={profile.linkedin} target="_blank" rel="noreferrer">
                      LinkedIn
                    </a>
                  )}
                </div>
              </div>
              <div className="portrait">
                {profile.photo_url ? (
                  <img src={asset(profile.photo_url)} alt={profile.name} />
                ) : (
                  <div className="ph">
                    Profile photo not set.
                    <br />
                    Upload one from the admin panel.
                  </div>
                )}
              </div>
            </section>

            {profile.strengths?.length > 0 && (
              <section>
                <h2>Core strengths</h2>
                <div className="grid-pills">
                  {profile.strengths.map((s) => (
                    <span className="pill" key={s}>
                      {s}
                    </span>
                  ))}
                </div>
              </section>
            )}

            {profile.skills?.length > 0 && (
              <section>
                <h2>Technical skills</h2>
                <div className="skill-grid">
                  {profile.skills.map((g) => (
                    <article className="skill-card" key={g.category}>
                      <h3>{g.category}</h3>
                      <ul>
                        {g.items.map((it) => (
                          <li key={it}>{it}</li>
                        ))}
                      </ul>
                    </article>
                  ))}
                </div>
              </section>
            )}

            <section id="projects">
              <h2>Selected projects</h2>
              {projects.map((p, i) => (
                <article className="project" key={p.id}>
                  <div className="project-head" onClick={() => setOpen(open === p.slug ? null : p.slug)}>
                    <div className="idx">{String(i + 1).padStart(2, "0")}</div>
                    <div>
                      <h3>
                        {p.name}
                        {p.subtitle ? ` — ${p.subtitle}` : ""}
                      </h3>
                      <p className="sub">{p.summary}</p>
                      <div className="facts">
                        {p.period && (
                          <span className="fact">
                            <b>Period</b> {p.period}
                          </span>
                        )}
                        {p.client && (
                          <span className="fact">
                            <b>Client</b> {p.client}
                          </span>
                        )}
                        {p.status && (
                          <span className="fact">
                            <b>Status</b> {p.status}
                          </span>
                        )}
                        {p.role && (
                          <span className="fact">
                            <b>Role</b> {p.role}
                          </span>
                        )}
                      </div>
                    </div>
                    <button className="ghost" type="button">
                      {open === p.slug ? "Close" : "Open"}
                    </button>
                  </div>
                  {open === p.slug && (
                    <div className="project-body">
                      {p.image_urls?.length > 0 && (
                        <div className="shots">
                          {p.image_urls.map((u) => (
                            <img key={u} src={asset(u)} alt={p.name} />
                          ))}
                        </div>
                      )}
                      {p.tech && p.tech.length > 0 && (
                        <div className="grid-pills" style={{ marginBottom: 16 }}>
                          {p.tech.map((t) => (
                            <span className="pill" key={t}>
                              {t}
                            </span>
                          ))}
                        </div>
                      )}
                      <Markdown text={p.content} />
                    </div>
                  )}
                </article>
              ))}
            </section>

            <section id="approach">
              <h2>How I work</h2>
              {profile.sections.map((s) => (
                <article className="section-block" key={s.slug}>
                  <h3 style={{ fontFamily: "var(--display)", fontSize: 26, margin: "0 0 8px" }}>{s.title}</h3>
                  <Markdown text={s.content} />
                </article>
              ))}
            </section>

            <section id="contact" className="footer">
              <div>
                <div className="kicker">Contact</div>
                <p>
                  {profile.email && (
                    <>
                      <a href={`mailto:${profile.email}`}>{profile.email}</a>
                      <br />
                    </>
                  )}
                  {profile.github && (
                    <>
                      <a href={profile.github}>{profile.github}</a>
                      <br />
                    </>
                  )}
                  {profile.linkedin && <a href={profile.linkedin}>{profile.linkedin}</a>}
                </p>
              </div>
              <div>Answers in the chat are limited to this portfolio.</div>
            </section>
          </>
        )}
      </main>
      <Chatbot />
    </>
  );
}
