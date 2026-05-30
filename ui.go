package main

import "html/template"

var dashTmpl = template.Must(template.New("dash").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="theme-color" content="#07080d">
<title>dispatch</title>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#07080d;--surface:rgba(255,255,255,0.04);--surface-hi:rgba(255,255,255,0.07);
  --border:rgba(255,255,255,0.08);--border-hi:rgba(255,255,255,0.14);
  --accent:#a78bfa;--green:#34d399;--red:#f87171;--amber:#fbbf24;
  --text:#e2e8f0;--text-dim:#94a3b8;--text-muted:#475569;
  --f:Inter,system-ui,-apple-system,sans-serif;
  --mono:'JetBrains Mono','SF Mono',monospace;
}
html,body{height:100%;background:var(--bg);color:var(--text);font-family:var(--f);
  font-size:14px;line-height:1.5;-webkit-font-smoothing:antialiased}

.layout{max-width:960px;margin:0 auto;padding:2rem 1rem 4rem}

.hdr{display:flex;align-items:center;justify-content:space-between;margin-bottom:2.5rem}
.logo{font-size:1.25rem;font-weight:700;letter-spacing:-0.02em}
.logo span{color:var(--accent)}
.hdr-right{display:flex;align-items:center;gap:0.75rem}
.badge{font-size:0.7rem;font-weight:600;letter-spacing:0.04em;padding:3px 10px;
  border-radius:999px;background:rgba(52,211,153,0.12);color:var(--green);
  border:1px solid rgba(52,211,153,0.2)}
.ver{font-size:0.72rem;color:var(--text-muted);font-family:var(--mono)}

.empty{display:flex;flex-direction:column;align-items:center;justify-content:center;
  gap:0.75rem;padding:4rem 2rem;text-align:center;
  border:1px dashed var(--border);border-radius:12px;color:var(--text-muted)}
.empty-ico{font-size:2rem;opacity:0.2}
.empty-txt{font-size:0.85rem;line-height:1.8;max-width:280px}
.empty code{font-family:var(--mono);font-size:0.8rem;background:var(--surface);
  padding:1px 6px;border-radius:4px;color:var(--text-dim)}

.section-label{font-size:0.7rem;font-weight:600;letter-spacing:0.08em;
  color:var(--text-muted);text-transform:uppercase;margin-bottom:0.75rem}

.stats{display:flex;gap:2rem;margin-bottom:1.75rem;padding:0.85rem 1.25rem;
  background:var(--surface);border:1px solid var(--border);border-radius:10px}
.stat{display:flex;flex-direction:column;gap:2px}
.stat-val{font-size:1.4rem;font-weight:700;letter-spacing:-0.02em}
.stat-val.green{color:var(--green)}
.stat-val.amber{color:var(--amber)}
.stat-lbl{font-size:0.68rem;color:var(--text-muted);letter-spacing:0.05em;text-transform:uppercase}

.workers{display:grid;gap:1rem;grid-template-columns:repeat(auto-fill,minmax(280px,1fr))}

.card{background:var(--surface);border:1px solid var(--border);
  border-radius:12px;overflow:hidden;transition:border-color 0.2s,background 0.2s}
.card:hover{background:var(--surface-hi);border-color:var(--border-hi)}

.card-head{display:flex;align-items:center;justify-content:space-between;
  padding:0.85rem 1rem;border-bottom:1px solid var(--border)}
.card-label{display:flex;align-items:center;gap:0.5rem;font-weight:600;font-size:0.88rem}
.dot{width:7px;height:7px;border-radius:50%;flex-shrink:0}
.dot.online{background:var(--green);box-shadow:0 0 6px rgba(52,211,153,0.5);
  animation:pdot 2.5s ease-in-out infinite}
.dot.offline{background:var(--text-muted)}
@keyframes pdot{0%,100%{opacity:1}50%{opacity:0.35}}
.card-status{font-size:0.7rem;color:var(--text-muted)}

.card-body{padding:0.75rem 1rem}

.session-list{display:flex;flex-direction:column;gap:4px}
.session-row{display:flex;align-items:center;gap:0.5rem;padding:5px 8px;
  border-radius:7px;background:rgba(255,255,255,0.025);transition:background 0.15s}
.session-row:hover{background:rgba(255,255,255,0.05)}
.session-dot{width:6px;height:6px;border-radius:50%;flex-shrink:0}
.session-dot.running{background:var(--green)}
.session-dot.stopped{background:var(--text-muted)}
.session-name{font-size:0.8rem;font-weight:500;font-family:var(--mono);flex:1;
  white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.session-cli{font-size:0.68rem;color:var(--text-muted);font-family:var(--mono);
  background:var(--surface);padding:1px 6px;border-radius:4px;flex-shrink:0}
.session-summary{font-size:0.72rem;color:var(--text-dim);margin-top:1px;
  padding-left:1rem;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.no-sessions{font-size:0.78rem;color:var(--text-muted);padding:2px 0}

.card-foot{display:flex;align-items:center;justify-content:space-between;
  padding:0.6rem 1rem;border-top:1px solid var(--border);margin-top:0.25rem}
.card-url{font-size:0.7rem;color:var(--text-muted);font-family:var(--mono);
  white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:55%}
.caps{display:flex;gap:4px;flex-wrap:wrap}
.cap{font-size:0.65rem;font-family:var(--mono);background:rgba(167,139,250,0.08);
  color:var(--accent);border:1px solid rgba(167,139,250,0.18);padding:1px 6px;border-radius:4px}
</style>
</head>
<body>
<div class="layout">

  <header class="hdr">
    <div class="logo">dis<span>patch</span></div>
    <div class="hdr-right">
      <span class="badge">hub</span>
      <span class="ver">v{{.Version}}</span>
    </div>
  </header>

  {{if .Workers}}

  <div class="stats">
    <div class="stat">
      <span class="stat-val green">{{.Online}}</span>
      <span class="stat-lbl">online</span>
    </div>
    <div class="stat">
      <span class="stat-val{{if .Offline}} amber{{end}}">{{.Offline}}</span>
      <span class="stat-lbl">offline</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{.Sessions}}</span>
      <span class="stat-lbl">sessions</span>
    </div>
  </div>

  <div class="section-label">Workers</div>
  <div class="workers">
  {{range .Workers}}
  <div class="card">
    <div class="card-head">
      <div class="card-label">
        <span class="dot {{if .Online}}online{{else}}offline{{end}}"></span>
        {{.Label}}
      </div>
      <span class="card-status">{{if .Online}}live{{else}}offline{{end}}</span>
    </div>
    <div class="card-body">
      {{if .Sessions}}
      <div class="session-list">
        {{range .Sessions}}
        <div>
          <div class="session-row">
            <span class="session-dot {{.Status}}"></span>
            <span class="session-name">{{.Name}}</span>
            <span class="session-cli">{{.CLI}}</span>
          </div>
          {{if .Summary}}<div class="session-summary">{{.Summary}}</div>{{end}}
        </div>
        {{end}}
      </div>
      {{else}}
      <div class="no-sessions">no sessions</div>
      {{end}}
    </div>
    <div class="card-foot">
      <span class="card-url">{{.URL}}</span>
      <div class="caps">
        {{range .Capabilities}}<span class="cap">{{.}}</span>{{end}}
      </div>
    </div>
  </div>
  {{end}}
  </div>

  {{else}}

  <div class="empty">
    <div class="empty-ico">◈</div>
    <div class="empty-txt">
      No workers registered yet.<br>
      Point a worker at this hub with:<br><br>
      <code>--hub http://&lt;this-host&gt;:8888</code>
    </div>
  </div>

  {{end}}

</div>
<script>
setTimeout(function(){ location.reload(); }, 15000);
</script>
</body>
</html>
`))
