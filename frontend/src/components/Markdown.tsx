import { Fragment, type ReactNode } from "react";

function inline(text: string) {
  const parts = text.split(/(\*\*[^*]+\*\*)/g);
  return parts.map((part, i) => {
    if (part.startsWith("**") && part.endsWith("**")) {
      return <strong key={i}>{part.slice(2, -2)}</strong>;
    }
    return <Fragment key={i}>{part}</Fragment>;
  });
}

export default function Markdown({ text }: { text: string }) {
  const blocks = text.replace(/\r\n/g, "\n").split(/\n```/);
  const nodes: ReactNode[] = [];
  blocks.forEach((chunk, i) => {
    if (i % 2 === 1) {
      const nl = chunk.indexOf("\n");
      const code = nl === -1 ? chunk : chunk.slice(nl + 1);
      nodes.push(<pre key={`c${i}`}>{code.replace(/\n$/, "")}</pre>);
      return;
    }
    chunk.split(/\n{2,}/).forEach((para, j) => {
      const line = para.trim();
      if (!line) return;
      if (line.startsWith("### ")) {
        nodes.push(<h4 key={`${i}-${j}`}>{inline(line.slice(4))}</h4>);
        return;
      }
      if (line.startsWith("## ") || line.startsWith("# ")) {
        nodes.push(<h3 key={`${i}-${j}`}>{inline(line.replace(/^#+ /, ""))}</h3>);
        return;
      }
      if (line.split("\n").every((l) => l.trim().startsWith("- "))) {
        nodes.push(
          <ul key={`${i}-${j}`}>
            {line.split("\n").map((l, k) => (
              <li key={k}>{inline(l.replace(/^- /, ""))}</li>
            ))}
          </ul>,
        );
        return;
      }
      nodes.push(<p key={`${i}-${j}`}>{inline(line)}</p>);
    });
  });
  return <div className="md">{nodes}</div>;
}
