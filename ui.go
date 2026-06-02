package main

import (
	"encoding/json"
	"html/template"
)

// dashTmpl is a pure JS-driven page — no Go template vars needed.
var dashTmpl = template.Must(template.New("dash").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,user-scalable=no">
<meta name="theme-color" content="#0d1117">
<meta name="apple-mobile-web-app-capable" content="yes">
<title>dispatch</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
  :root {
    --bg:#0d1117; --surface:rgba(255,255,255,0.04); --border:rgba(255,255,255,0.1);
    --border-hi:rgba(255,255,255,0.2); --accent:#58a6ff; --accent-dim:rgba(88,166,255,0.15);
    --green:#3fb950; --green-dim:rgba(63,185,80,0.15); --red:#f85149;
    --red-dim:rgba(248,81,73,0.15); --amber:#d29922;
    --text:#e2e8f0; --text-secondary:#8b949e;
    --mono:'JetBrains Mono','SF Mono','Consolas',monospace;
  }
  *{box-sizing:border-box;margin:0;padding:0}
  body{font-family:var(--mono);background:var(--bg);color:var(--text);
    padding:1.25rem 1rem 6rem;-webkit-text-size-adjust:100%;min-height:100vh}

  .title{display:flex;align-items:center;justify-content:space-between;
    font-size:1rem;font-weight:700;margin-bottom:2rem;
    padding-bottom:1rem;border-bottom:1px solid var(--border)}

  .machine-tabs{display:flex;align-items:center;gap:.3rem;flex-wrap:wrap}
  .machine-tab{display:flex;align-items:center;gap:.3rem;padding:.2rem .6rem;
    border-radius:20px;font-size:.75rem;font-family:var(--mono);font-weight:500;
    border:1px solid var(--border);background:none;color:var(--text-secondary);
    cursor:pointer;transition:all .15s;white-space:nowrap}
  .machine-tab:hover{border-color:var(--accent);color:var(--accent)}
  .machine-tab.active{background:var(--accent);border-color:var(--accent);color:#0d1117}
  .machine-tab.active .status-dot{background:rgba(0,0,0,.4)}
  .status-dot{width:6px;height:6px;border-radius:50%;flex-shrink:0;background:var(--border)}
  .status-dot.online{background:var(--green)}
  .status-dot.offline{background:var(--red)}

  .toolbar{display:flex;align-items:center;margin-bottom:2rem}
  .new-btn{margin-left:auto;padding:.5rem 1.1rem;border-radius:8px;
    font-size:.85rem;font-weight:700;font-family:var(--mono);cursor:pointer;
    background:var(--accent-dim);color:var(--accent);
    border:1px solid rgba(88,166,255,.3);transition:background .15s}
  .new-btn:hover{background:rgba(88,166,255,.25)}

  .section-label{font-size:.75rem;color:var(--text-secondary);text-transform:uppercase;
    letter-spacing:.1em;margin-bottom:.75rem;font-weight:500}

  .card{background:var(--surface);border:1px solid var(--border);border-radius:8px;
    margin-bottom:.75rem;transition:border-color .15s;position:relative;overflow:visible;
    animation:cardIn .4s cubic-bezier(.16,1,.3,1) both}
  .card:hover{border-color:var(--border-hi)}
  @keyframes cardIn{from{opacity:0;transform:translateY(8px)}to{opacity:1;transform:translateY(0)}}

  .card-body{padding:.9rem 1rem .75rem;cursor:pointer;border-radius:8px 8px 0 0;transition:background .12s}
  .card-body:hover{background:rgba(255,255,255,.03)}
  .card-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:.3rem}
  .card-name{font-weight:700;font-size:.95rem;letter-spacing:.01em}
  .card-cli{font-size:.75rem;color:var(--text-secondary)}
  .card-summary{font-size:.75rem;color:var(--text-secondary);margin-top:.35rem}

  .badge{font-size:.72rem;padding:.18rem .6rem;border-radius:9999px;font-weight:600;
    letter-spacing:.02em;display:inline-flex;align-items:center;gap:.3rem}
  .badge-running{background:var(--green-dim);color:var(--green)}
  .badge-running::before{content:'';width:5px;height:5px;background:var(--green);
    border-radius:50%;animation:pulse 2s ease-in-out infinite}
  @keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}
  .badge-stopped{background:rgba(255,255,255,.06);color:var(--text-secondary)}

  .card-footer{display:flex;justify-content:flex-end;align-items:center;
    padding:.25rem .6rem;border-top:1px solid var(--border);position:relative}
  .more-btn{background:none;border:1px solid transparent;color:var(--text-secondary);
    font-family:var(--mono);font-size:.72rem;font-weight:600;letter-spacing:.04em;
    cursor:pointer;padding:.2rem .5rem;border-radius:4px;transition:all .12s}
  .more-btn:hover{background:rgba(255,255,255,.06);border-color:var(--border);color:var(--text)}

  .card-menu{position:absolute;bottom:calc(100% + 4px);right:0;min-width:185px;
    background:#1c2128;border:1px solid var(--border-hi);border-radius:8px;
    padding:.3rem;z-index:50;box-shadow:0 8px 24px rgba(0,0,0,.6);display:none}
  .card-menu.open{display:block}
  .menu-item{display:block;width:100%;padding:.5rem .7rem;border-radius:5px;
    font-size:.8rem;font-weight:600;font-family:var(--mono);cursor:pointer;
    border:none;background:none;text-align:left;transition:background .1s;color:var(--text)}
  .menu-item:hover{background:rgba(255,255,255,.06)}
  .menu-item.kill{color:var(--red)}
  .menu-item.resume{color:var(--green)}
  .menu-item.restart{color:var(--amber)}
  .menu-divider{height:1px;background:var(--border);margin:.25rem .3rem}
  .menu-path{padding:.4rem .7rem;font-size:.7rem;color:var(--text-secondary);
    overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

  .worker-info{display:flex;gap:1rem;flex-wrap:wrap;margin-bottom:1rem;
    font-size:.72rem;color:var(--text-secondary)}
  .worker-info a{color:var(--accent);text-decoration:none;font-family:var(--mono)}
  .worker-info a:hover{text-decoration:underline}
  .empty{color:var(--text-secondary);text-align:center;padding:3rem 0;
    font-size:.85rem;border:1px dashed var(--border);border-radius:6px}

  .overlay{position:fixed;inset:0;background:rgba(0,0,0,.7);
    backdrop-filter:blur(10px);z-index:100;display:flex;align-items:center;
    justify-content:center;opacity:0;pointer-events:none;transition:opacity .2s}
  .overlay.show{opacity:1;pointer-events:auto}
  .modal{background:#161b22;border:1px solid var(--border-hi);border-radius:8px;
    width:100%;max-width:400px;overflow:hidden;
    transform:translateY(16px);transition:transform .2s cubic-bezier(.16,1,.3,1)}
  .overlay.show .modal{transform:translateY(0)}
  .modal-head{display:flex;align-items:center;justify-content:space-between;
    padding:.85rem 1.1rem;border-bottom:1px solid var(--border)}
  .modal-title{font-size:.9rem;font-weight:700}
  .modal-close{background:none;border:none;color:var(--text-secondary);
    font-size:1.2rem;cursor:pointer;padding:2px 4px;transition:color .12s}
  .modal-close:hover{color:var(--text)}
  .modal-body{padding:1.1rem;display:flex;flex-direction:column;gap:.85rem}
  .fg{display:flex;flex-direction:column;gap:.35rem}
  .fg label{font-size:.7rem;font-weight:600;color:var(--text-secondary);
    text-transform:uppercase;letter-spacing:.06em}
  .fg select,.fg input{background:rgba(255,255,255,.06);border:1px solid var(--border);
    color:var(--text);font-family:var(--mono);font-size:.85rem;
    padding:.45rem .65rem;border-radius:6px;outline:none;transition:border-color .15s;width:100%}
  .fg select:focus,.fg input:focus{border-color:rgba(88,166,255,.5)}
  .fg select{cursor:pointer;appearance:none}
  .modal-foot{padding:.75rem 1.1rem;border-top:1px solid var(--border);
    display:flex;justify-content:flex-end;gap:.5rem}
  .btn{font-size:.82rem;font-weight:700;padding:.45rem 1rem;font-family:var(--mono);
    border-radius:6px;border:none;cursor:pointer;transition:all .12s}
  .btn-cancel{background:rgba(255,255,255,.06);color:var(--text-secondary);border:1px solid var(--border)}
  .btn-cancel:hover{color:var(--text)}
  .btn-ok{background:var(--accent-dim);color:var(--accent);border:1px solid rgba(88,166,255,.3)}
  .btn-ok:hover:not(:disabled){background:rgba(88,166,255,.25)}
  .btn-ok:disabled{opacity:.4;cursor:not-allowed}
</style>
</head>
<body>

<div class="title">
  <div class="machine-tabs" id="machine-tabs">
    <button class="machine-tab active" id="tab-all" onclick="setTab('all')">All</button>
  </div>
</div>

<div class="toolbar">
  <div class="section-label">Sessions</div>
  <button class="new-btn" onclick="openSpawn()">+ New session</button>
</div>

<div id="worker-info"></div>
<div id="instances"></div>

<div class="overlay" id="overlay" onclick="if(event.target===this)closeModal()">
  <div class="modal">
    <div class="modal-head">
      <span class="modal-title">New session</span>
      <button class="modal-close" onclick="closeModal()">&#x2715;</button>
    </div>
    <div class="modal-body">
      <div class="fg" id="m-machine-fg">
        <label>Machine</label>
        <select id="m-machine" onchange="updateCaps()"></select>
      </div>
      <div class="fg">
        <label>Tool</label>
        <select id="m-cli"></select>
      </div>
      <div class="fg">
        <label>Directory</label>
        <input id="m-dir" type="text" placeholder="~">
      </div>
      <div class="fg">
        <label>Session name (optional)</label>
        <input id="m-name" type="text" placeholder="auto-generated">
      </div>
    </div>
    <div class="modal-foot">
      <button class="btn btn-cancel" onclick="closeModal()">Cancel</button>
      <button class="btn btn-ok" id="m-submit" onclick="submitSpawn()">Spawn</button>
    </div>
  </div>
</div>

<script>
var lastWorkers = [];
var activeTab = sessionStorage.getItem('activeTab') || 'all';

function esc(s) {
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function setTab(id) {
  activeTab = id;
  sessionStorage.setItem('activeTab', id);
  document.querySelectorAll('.machine-tab').forEach(function(t){ t.classList.remove('active'); });
  var tab = document.getElementById('tab-' + id);
  if (tab) tab.classList.add('active');
  renderSessions();
}

function updateTabs(workers) {
  var tabs = document.getElementById('machine-tabs');
  tabs.querySelectorAll('.machine-tab:not(#tab-all)').forEach(function(t){ t.remove(); });
  workers.forEach(function(w) {
    var btn = document.createElement('button');
    btn.id = 'tab-' + w.id;
    btn.className = 'machine-tab' + (activeTab === w.id ? ' active' : '');
    var dot = '<span class="status-dot ' + (w.online ? 'online' : 'offline') + '"></span>';
    btn.innerHTML = dot + esc(w.label);
    btn.onclick = function(){ setTab(w.id); };
    tabs.appendChild(btn);
  });
  document.getElementById('tab-all').className = 'machine-tab' + (activeTab === 'all' ? ' active' : '');
}

function renderSessions() {
  var el = document.getElementById('instances');
  var items = [];
  if (activeTab === 'all') {
    lastWorkers.forEach(function(w) {
      (w.sessions || []).forEach(function(s) {
        items.push({s:s, wid:w.id, wlabel:w.label, wurl:w.url});
      });
    });
  } else {
    var w = lastWorkers.find(function(wk){ return wk.id === activeTab; });
    if (w) (w.sessions || []).forEach(function(s){ items.push({s:s, wid:w.id, wlabel:w.label, wurl:w.url}); });
  }

  document.getElementById('worker-info').innerHTML = '';

  if (!items.length) {
    el.innerHTML = '<div class="empty">No sessions</div>';
    return;
  }

  var html = '';
  items.forEach(function(item, i) {
    var s = item.s, wid = item.wid;
    var running = s.status === 'running';
    var cliText = activeTab === 'all' ? (esc(s.cli||'terminal') + ' · ' + esc(item.wlabel)) : esc(s.cli||'terminal');
    var mid = esc(s.name) + '-' + esc(wid);
    html += '<div class="card" style="animation-delay:' + (i*0.05) + 's">';
    html += '<div class="card-body" onclick="openSession(\'' + esc(wid) + '\',\'' + esc(s.name) + '\')">';
    html += '<div class="card-header">';
    html += '<span class="card-name">' + esc(s.name) + '</span>';
    html += '<span class="badge badge-' + s.status + '">' + s.status + '</span>';
    html += '</div>';
    html += '<div class="card-cli">' + cliText + '</div>';
    if (s.summary) html += '<div class="card-summary">' + esc(s.summary) + '</div>';
    html += '</div>';
    html += '<div class="card-footer">';
    html += '<button class="more-btn" onclick="toggleMenu(\'' + mid + '\',event)">More</button>';
    html += '<div class="card-menu" id="menu-' + mid + '">';
    if (running) {
      html += '<button class="menu-item kill" onclick="doAction(\'kill\',\'' + esc(wid) + '\',\'' + esc(s.name) + '\')">Kill session</button>';
      html += '<button class="menu-item restart" onclick="doAction(\'restart\',\'' + esc(wid) + '\',\'' + esc(s.name) + '\')">Restart</button>';
    } else {
      html += '<button class="menu-item resume" onclick="doAction(\'resume\',\'' + esc(wid) + '\',\'' + esc(s.name) + '\')">Resume</button>';
      html += '<button class="menu-item kill" onclick="doAction(\'delete\',\'' + esc(wid) + '\',\'' + esc(s.name) + '\')">Delete</button>';
    }
    if (s.dir) { html += '<div class="menu-divider"></div><div class="menu-path">' + esc(s.dir) + '</div>'; }
    html += '</div></div></div>';
  });
  el.innerHTML = html;
}

function render(workers) {
  lastWorkers = workers;
  updateTabs(workers);
  renderSessions();
}

// fetchWorkerSessions fetches live session state directly from a worker's API.
// Falls back to the heartbeat-cached sessions if the worker is unreachable.
function fetchWorkerSessions(w) {
  if (!w.online) return Promise.resolve(w);
  return fetch('/api/workers/' + encodeURIComponent(w.id) + '/instances')
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(sessions) {
      if (!sessions) return w;
      return Object.assign({}, w, {sessions: sessions});
    })
    .catch(function() { return w; });
}

// load fetches worker metadata (online/offline) then fetches live sessions
// from each online worker's API in parallel. Heartbeat data is only used
// for online/offline status and worker metadata, never for session state.
function load() {
  fetch('/api/workers')
    .then(function(r) { return r.json(); })
    .then(function(workers) {
      if (!workers || !workers.length) { render([]); return; }
      return Promise.all(workers.map(fetchWorkerSessions));
    })
    .then(function(workers) { if (workers) render(workers); })
    .catch(function(){});
}

// refreshWorker re-fetches live sessions for a single worker and re-renders.
function refreshWorker(wid) {
  var w = lastWorkers.find(function(wk) { return wk.id === wid; });
  if (!w) return;
  fetchWorkerSessions(w).then(function(updated) {
    lastWorkers = lastWorkers.map(function(wk) { return wk.id === wid ? updated : wk; });
    renderSessions();
  });
}

function openSession(wid, name) {
  location.href = '/session/' + encodeURIComponent(wid) + '/' + encodeURIComponent(name);
}

function doAction(action, wid, name) {
  closeAllMenus();
  fetch('/api/workers/' + encodeURIComponent(wid) + '/' + action + '/' + encodeURIComponent(name), {method:'POST'})
    .then(function(r) {
      return r.json().then(function(d) {
        if (!r.ok) {
          // Kill on already-stopped: not an error, just refresh to get real state.
          if (action === 'kill' && d.message && d.message.indexOf('not running') !== -1) {
            refreshWorker(wid); return;
          }
          alert(d.error || d.message || 'Action failed'); return;
        }
        // Optimistic update so the UI responds instantly, then fetch real state.
        if (action === 'delete') {
          lastWorkers.forEach(function(w) {
            if (w.id === wid) w.sessions = (w.sessions||[]).filter(function(s){ return s.name !== name; });
          });
          renderSessions();
        } else {
          lastWorkers.forEach(function(w) {
            if (w.id !== wid) return;
            (w.sessions||[]).forEach(function(s) {
              if (s.name !== name) return;
              if (action === 'kill')    s.status = 'stopped';
              if (action === 'resume')  s.status = 'running';
              if (action === 'restart') s.status = 'running';
            });
          });
          renderSessions();
        }
        // Fetch real state from worker after a short delay for the action to land.
        setTimeout(function() { refreshWorker(wid); }, 800);
      });
    }).catch(function(){});
}

function toggleMenu(mid, e) {
  e.stopPropagation();
  var menu = document.getElementById('menu-' + mid);
  var open = menu.classList.contains('open');
  closeAllMenus();
  if (!open) menu.classList.add('open');
}
function closeAllMenus() {
  document.querySelectorAll('.card-menu.open').forEach(function(m){ m.classList.remove('open'); });
}
document.addEventListener('click', closeAllMenus);

function openSpawn() {
  var online = lastWorkers.filter(function(w){ return w.online; });
  if (!online.length) { alert('No online workers'); return; }
  var sel = document.getElementById('m-machine');
  sel.innerHTML = '';
  online.forEach(function(w) {
    var opt = document.createElement('option');
    opt.value = w.id; opt.textContent = w.label;
    if (w.id === activeTab) opt.selected = true;
    sel.appendChild(opt);
  });
  updateCaps();
  document.getElementById('m-dir').value = '';
  document.getElementById('m-name').value = '';
  document.getElementById('overlay').classList.add('show');
  document.getElementById('m-dir').focus();
}

function updateCaps() {
  var wid = document.getElementById('m-machine').value;
  var w = lastWorkers.find(function(wk){ return wk.id === wid; });
  var caps = (w && w.capabilities) ? w.capabilities : ['terminal'];
  var sel = document.getElementById('m-cli');
  sel.innerHTML = '';
  caps.forEach(function(c){ sel.innerHTML += '<option value="' + esc(c) + '">' + esc(c) + '</option>'; });
}

function closeModal() { document.getElementById('overlay').classList.remove('show'); }

function submitSpawn() {
  var btn = document.getElementById('m-submit');
  btn.disabled = true; btn.textContent = 'Spawning...';
  var wid = document.getElementById('m-machine').value;
  var payload = {
    cli: document.getElementById('m-cli').value,
    dir: document.getElementById('m-dir').value || '~',
    name: document.getElementById('m-name').value || ''
  };
  fetch('/api/workers/' + encodeURIComponent(wid) + '/spawn', {
    method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(payload)
  })
  .then(function(r){ return r.json(); })
  .then(function(d) {
    btn.disabled = false; btn.textContent = 'Spawn';
    if (d.error) { alert(d.error); return; }
    closeModal();
    var name = d.name || payload.name;
    if (name) location.href = '/session/' + encodeURIComponent(wid) + '/' + encodeURIComponent(name);
    else load();
  })
  .catch(function(){ btn.disabled = false; btn.textContent = 'Spawn'; alert('Spawn failed'); });
}

document.addEventListener('keydown', function(e){ if (e.key === 'Escape') closeModal(); });

load();
setInterval(load, 30000);
</script>
</body>
</html>
`))

// sessionTmpl renders a terminal session page that connects directly to the worker WebSocket.
var sessionTmpl = template.Must(template.New("session").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,user-scalable=no,viewport-fit=cover">
<meta name="theme-color" content="#07080d">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
<title>{{.SessionName}} — dispatch</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/css/xterm.min.css">
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0d1117;--surface:rgba(255,255,255,0.04);--border:rgba(255,255,255,0.1);
  --accent:#58a6ff;--green:#3fb950;--green-dim:rgba(63,185,80,0.15);
  --red:#f85149;--red-dim:rgba(248,81,73,0.15);--amber:#d29922;
  --text:#e2e8f0;--text-dim:#94a3b8;--text-muted:#8b949e;
  --mono:'JetBrains Mono','SF Mono',monospace;
  --sat:env(safe-area-inset-top,0px);--sab:env(safe-area-inset-bottom,0px);
  --sal:env(safe-area-inset-left,0px);--sar:env(safe-area-inset-right,0px);
}
html,body{height:100%;overflow:hidden;background:var(--bg);color:var(--text);
  font-family:var(--mono);-webkit-font-smoothing:antialiased;-webkit-text-size-adjust:100%}
.layout{display:flex;flex-direction:column;height:100%;height:100dvh}

/* topbar */
.topbar{flex-shrink:0;display:flex;align-items:center;justify-content:space-between;
  padding:calc(var(--sat) + 8px) calc(var(--sar) + 12px) 8px calc(var(--sal) + 6px);
  background:rgba(13,17,23,.9);backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);
  border-bottom:1px solid var(--border);z-index:10}
.topbar-left{display:flex;align-items:center;gap:.5rem;flex:1;min-width:0;overflow:hidden}
.topbar-right{display:flex;align-items:center;gap:.4rem;flex-shrink:0}
.back{color:var(--accent);background:transparent;border:none;
  font-size:1.5rem;cursor:pointer;padding:.25rem .6rem .25rem 0;
  display:flex;align-items:center;line-height:1;text-decoration:none;
  transition:opacity .15s;flex-shrink:0}
.back:hover{opacity:.7}
.session-label{font-size:.85rem;font-weight:700;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.session-sep{color:var(--text-muted);flex-shrink:0}
.session-sub{font-size:.85rem;color:var(--text-muted);white-space:nowrap;flex-shrink:0}
.badge{font-size:.66rem;font-weight:600;letter-spacing:.04em;
  padding:2px 8px;border-radius:999px;flex-shrink:0}
.badge-live{background:var(--green-dim);color:var(--green);border:1px solid rgba(52,211,153,.25)}
.badge-stopped{background:rgba(255,255,255,.04);color:var(--text-muted);border:1px solid var(--border)}
.badge-connecting{background:rgba(167,139,250,.1);color:var(--accent);border:1px solid rgba(167,139,250,.25)}

/* menu */
.menu-wrap{position:relative;flex-shrink:0}
.menu-btn{background:transparent;border:none;color:var(--text-dim);
  cursor:pointer;font-size:1.2rem;padding:4px 6px;border-radius:6px;
  transition:background .15s;line-height:1}
.menu-btn:hover{background:var(--surface)}
.dropdown{position:absolute;right:0;top:calc(100% + 4px);
  background:#0f1018;border:1px solid var(--border);border-radius:10px;
  min-width:140px;overflow:hidden;z-index:50;
  box-shadow:0 8px 24px rgba(0,0,0,.5);display:none}
.dropdown.open{display:block}
.ditem{display:block;width:100%;text-align:left;background:transparent;border:none;
  color:var(--text-dim);font-family:var(--mono);font-size:.8rem;
  padding:.5rem .85rem;cursor:pointer;transition:background .1s,color .1s}
.ditem:hover{background:var(--surface);color:var(--text)}
.ditem.red:hover{color:var(--red)}
.dsep{height:1px;background:var(--border);margin:2px 0}

/* terminal */
#term{flex:1;overflow:hidden}
#term .xterm{height:100%}
#term .xterm-viewport{overflow-y:auto!important}

/* action bar */
.abar{flex-shrink:0;display:none;
  padding:6px calc(var(--sar) + 8px) calc(var(--sab) + 6px) calc(var(--sal) + 8px);
  background:#0a0e14;border-top:1px solid var(--border);
  flex-wrap:nowrap;gap:4px;overflow-x:auto;align-items:center}
.abar::-webkit-scrollbar{display:none}
.abar.visible{display:flex}
.ak{height:30px;min-width:40px;padding:0 8px;border-radius:6px;
  background:var(--surface);border:1px solid var(--border);
  color:var(--text-dim);font-family:var(--mono);font-size:.72rem;font-weight:500;
  cursor:pointer;display:flex;align-items:center;justify-content:center;
  transition:background .1s;-webkit-tap-highlight-color:transparent;touch-action:manipulation;flex-shrink:0}
.ak:active{background:rgba(255,255,255,.1)}
.ak.enter{background:rgba(88,166,255,.1);border-color:rgba(88,166,255,.3);color:var(--accent)}
.ak.ctrlc{background:var(--red-dim);border-color:rgba(248,113,113,.25);color:var(--red)}
.ak-sep{width:1px;height:16px;background:var(--border);flex-shrink:0;margin:0 2px}
</style>
</head>
<body>
<div class="layout">

  <div class="topbar">
    <div class="topbar-left">
      <a class="back" href="/" title="Back">&#8592;</a>
      <span class="session-label">{{.SessionName}}</span>
      <span class="session-sep">·</span>
      <span class="session-sub">{{.WorkerLabel}}</span>
      <span class="badge badge-connecting" id="badge">connecting</span>
    </div>
    <div class="topbar-right">
      <div class="menu-wrap">
        <button class="menu-btn" id="menu-btn" onclick="toggleMenu()">&#8943;</button>
        <div class="dropdown" id="dropdown">
          <button class="ditem" id="dd-kill"    onclick="sessionAction('kill');closeMenu()"    style="display:{{if eq .SessionStatus "running"}}block{{else}}none{{end}}">Kill session</button>
          <button class="ditem" id="dd-restart" onclick="sessionAction('restart');closeMenu()" style="display:{{if eq .SessionStatus "running"}}block{{else}}none{{end}}">Restart</button>
          <button class="ditem" id="dd-resume"  onclick="sessionAction('resume');closeMenu()"  style="display:{{if eq .SessionStatus "running"}}none{{else}}block{{end}}">Resume</button>
          <button class="ditem red" id="dd-delete" onclick="sessionAction('delete');closeMenu()" style="display:{{if eq .SessionStatus "running"}}none{{else}}block{{end}}">Delete</button>
          <div class="dsep"></div>
          <button class="ditem" onclick="sendCtrlC();closeMenu()">Interrupt (^C)</button>
          <button class="ditem" onclick="document.getElementById('file-input').click();closeMenu()">Attach file</button>
          <div class="dsep"></div>
          <a class="ditem" href="{{.WorkerURL}}/desktop" target="_blank" onclick="closeMenu()" style="text-decoration:none;display:block">Desktop (VNC)</a>
        </div>
      </div>
    </div>
  </div>

  <input type="file" id="file-input" accept="image/*,application/pdf" multiple style="display:none" onchange="handleFileInput(this)">
  <div id="term"></div>

  <div class="abar" id="abar">
    <button class="ak" onclick="sendKey('\x1b[A')">&#8593;</button>
    <button class="ak" onclick="sendKey('\x1b[B')">&#8595;</button>
    <button class="ak" onclick="sendKey('\x1b[D')">&#8592;</button>
    <button class="ak" onclick="sendKey('\x1b[C')">&#8594;</button>
    <div class="ak-sep"></div>
    <button class="ak" onclick="sendKey('\x1b')">esc</button>
    <button class="ak" onclick="sendKey('\t')">tab</button>
    <button class="ak enter" onclick="sendKey('\r')">&#8629;</button>
    <div class="ak-sep"></div>
    <button class="ak ctrlc" onclick="sendCtrlC()">^C</button>
  </div>

</div>

<script src="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/lib/xterm.js"></script>
<script src="https://cdn.jsdelivr.net/npm/@xterm/addon-fit@0.11.0/lib/addon-fit.js"></script>
<script>
var WS_URL      = {{.WSURL}};
var WS_TOKEN    = {{.WorkerToken}};
var WORKER_ID   = {{.WorkerID | js}};
var WORKER_URL  = {{.WorkerURL | js}};
var SESS_NAME   = {{.SessionName | js}};
var INIT_STATUS = {{.SessionStatus | js}};

var term = new Terminal({
  cursorBlink:true, cursorStyle:'bar',
  fontSize:13, fontFamily:"'JetBrains Mono','SF Mono','Menlo','Consolas',monospace",
  fontWeight:'400', fontWeightBold:'700',
  theme:{
    background:'#07080d', foreground:'#e2e8f0', cursor:'#a78bfa',
    selectionBackground:'rgba(167,139,250,0.25)',
    black:'#1a1a2e', red:'#f87171', green:'#34d399', yellow:'#fbbf24',
    blue:'#60a5fa', magenta:'#a78bfa', cyan:'#34d4d4', white:'#e2e8f0',
    brightBlack:'#475569', brightRed:'#fca5a5', brightGreen:'#6ee7b7',
    brightYellow:'#fde68a', brightBlue:'#93c5fd', brightMagenta:'#c4b5fd',
    brightCyan:'#67e8f9', brightWhite:'#f8fafc',
  },
  scrollback:5000, convertEol:false,
});

var fit = new FitAddon.FitAddon();
term.loadAddon(fit);
term.open(document.getElementById('term'));
fit.fit();

window.addEventListener('resize', function(){ fit.fit(); sendResize(); });
new ResizeObserver(function(){ fit.fit(); sendResize(); }).observe(document.getElementById('term'));

var ws = null;
var badge = document.getElementById('badge');

function setBadge(state) {
  badge.className = 'badge badge-' + state;
  badge.textContent = state;
}

function sendResize() {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({type:'resize', cols:term.cols, rows:term.rows}));
  }
}

function sendKey(seq) {
  term.scrollToBottom();
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({type:'input', data:seq}));
  }
  term.focus();
}

function sendCtrlC() { sendKey('\x03'); }

term.onData(function(data) {
  term.scrollToBottom();
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({type:'input', data:data}));
  }
});

// touch overlay for mobile scrolling
if (navigator.maxTouchPoints > 0) {
  document.getElementById('abar').classList.add('visible');
  var tc = document.getElementById('term');
  var ov = document.createElement('div');
  ov.style.cssText = 'position:absolute;inset:0;z-index:10;touch-action:none;-webkit-tap-highlight-color:transparent;';
  tc.style.position = 'relative';
  tc.appendChild(ov);
  var ty0=0,tyL=0,tL=0,vel=0,acc=0,moved=false,raf=null;
  function cancelM(){ if(raf){cancelAnimationFrame(raf);raf=null;} }
  function runM(){
    if(Math.abs(vel)<.15){vel=0;return;}
    acc+=vel; var ln=Math.trunc(acc/16);
    if(ln){term.scrollLines(ln);acc-=ln*16;}
    vel*=.94; raf=requestAnimationFrame(runM);
  }
  ov.addEventListener('touchstart',function(e){
    e.preventDefault(); if(e.touches.length>1){moved=true;return;}
    cancelM(); ty0=tyL=e.touches[0].clientY; tL=Date.now(); vel=0;acc=0;moved=false;
  },{passive:false});
  ov.addEventListener('touchmove',function(e){
    e.preventDefault(); if(e.touches.length!==1)return;
    var y=e.touches[0].clientY, dy=tyL-y, now=Date.now();
    vel=dy/Math.max(now-tL,1)*16; acc+=dy;
    var ln=Math.trunc(acc/16); if(ln){term.scrollLines(ln);acc-=ln*16;}
    if(Math.abs(y-ty0)>5)moved=true; tyL=y;tL=now;
  },{passive:false});
  ov.addEventListener('touchend',function(e){
    e.preventDefault();
    if(moved)raf=requestAnimationFrame(runM); else term.focus();
  },{passive:false});
  // On desktop browsers that report maxTouchPoints>0, the overlay intercepts
  // mouse clicks and prevents xterm from receiving focus. Forward clicks to focus.
  ov.addEventListener('click',function(){ term.focus(); });
}

function connect() {
  setBadge(INIT_STATUS === 'running' ? 'connecting' : 'stopped');
  ws = new WebSocket(WS_URL);
  ws.binaryType = 'arraybuffer';
  ws.onopen = function() {
    if (WS_TOKEN) ws.send(JSON.stringify({type: 'auth', token: WS_TOKEN}));
    setBadge('live');
    sendResize();
    term.focus();
    // Fetch real session status now that the WS is up — avoids racing with the
    // initial fetchSessionStatus() call and ensures buttons always reflect truth.
    fetchSessionStatus();
  };
  ws.onmessage = function(e) {
    var data = new Uint8Array(e.data);
    var buf = term.buffer.active;
    var atBot = buf.viewportY + term.rows >= buf.length - 3;
    term.write(data, function(){ if(atBot) term.scrollToBottom(); });
  };
  ws.onclose = function(e) {
    ws = null;
    // Code 4000 = worker confirmed session is stopped. Show stopped and do not reconnect.
    if (e.code === 4000) { setBadge('stopped'); return; }
    // Stopped sessions: show output once then stop reconnecting.
    if (INIT_STATUS !== 'running') { setBadge('stopped'); return; }
    setBadge('connecting');
    setTimeout(connect, 2000);
  };
  ws.onerror = function() { setBadge('connecting'); };
}

connect();
fetchSessionStatus();

// updateSessionButtons shows/hides the kill/restart vs resume/delete buttons
// based on actual session status. Called on load and after every action.
function updateSessionButtons(status) {
  var running = status === 'running';
  document.getElementById('dd-kill').style.display    = running ? 'block' : 'none';
  document.getElementById('dd-restart').style.display = running ? 'block' : 'none';
  document.getElementById('dd-resume').style.display  = running ? 'none'  : 'block';
  document.getElementById('dd-delete').style.display  = running ? 'none'  : 'block';
}

// fetchSessionStatus fetches live session state from the worker API and
// updates the dropdown buttons. Called on page load so the buttons always
// reflect actual state, not the stale status baked into the template.
function fetchSessionStatus() {
  fetch('/api/workers/' + encodeURIComponent(WORKER_ID) + '/instances')
    .then(function(r) { return r.json(); })
    .then(function(sessions) {
      var sess = sessions.find(function(s) { return s.name === SESS_NAME; });
      if (!sess) return;
      INIT_STATUS = sess.status;
      updateSessionButtons(sess.status);
      // Correct the badge if the worker says stopped but WS shows live.
      if (sess.status !== 'running') setBadge('stopped');
    })
    .catch(function(){});
}

function sessionAction(action) {
  fetch('/api/workers/' + encodeURIComponent(WORKER_ID) + '/' + action + '/' + encodeURIComponent(SESS_NAME), {method:'POST'})
    .then(function(r) {
      return r.json().then(function(d) {
        if (!r.ok) { alert(d.error || d.message || action + ' failed'); return; }
        if (action === 'delete') { location.href = '/'; return; }
        if (action === 'kill') {
          INIT_STATUS = 'stopped';
          updateSessionButtons('stopped');
          setBadge('stopped');
        } else {
          // resume or restart: reconnect WS so the terminal becomes live again.
          INIT_STATUS = 'running';
          updateSessionButtons('running');
          setBadge('connecting');
          if (ws) { ws.onclose = null; ws.close(); ws = null; }
          setTimeout(connect, 500);
        }
      });
    })
    .catch(function(){});
}

function toggleMenu() {
  document.getElementById('dropdown').classList.toggle('open');
}
function closeMenu() {
  document.getElementById('dropdown').classList.remove('open');
}
document.addEventListener('click', function(e) {
  var btn = document.getElementById('menu-btn');
  var dd = document.getElementById('dropdown');
  if (!btn.contains(e.target) && !dd.contains(e.target)) dd.classList.remove('open');
});

// File upload — sends directly to worker with worker token auth
async function uploadFile(file) {
  if (!file) return null;
  var form = new FormData();
  form.append('file', file);
  try {
    var res = await fetch(WORKER_URL + '/upload/' + encodeURIComponent(SESS_NAME), {
      method: 'POST',
      headers: {'Authorization': 'Bearer ' + WS_TOKEN},
      body: form
    });
    var d = await res.json();
    return d.status === 'ok' ? d.path : null;
  } catch(e) { return null; }
}

async function handleFileInput(input) {
  var files = Array.from(input.files || []);
  input.value = '';
  if (!files.length) return;
  var paths = (await Promise.all(files.map(uploadFile))).filter(Boolean);
  if (!paths.length) return;
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({type:'input', data: paths.join(' ')}));
  }
  term.focus();
}

// Paste image from clipboard
document.addEventListener('paste', function(e) {
  var items = e.clipboardData && e.clipboardData.items;
  if (!items) return;
  for (var i = 0; i < items.length; i++) {
    if (items[i].type.startsWith('image/')) {
      e.preventDefault();
      uploadFile(items[i].getAsFile()).then(function(path) {
        if (path && ws && ws.readyState === WebSocket.OPEN)
          ws.send(JSON.stringify({type:'input', data: path}));
      });
      return;
    }
  }
});
</script>
</body>
</html>
`))

// sessionData is passed to sessionTmpl.
type sessionData struct {
	WorkerID      string
	WorkerLabel   string
	WorkerURL     string
	SessionName   string
	SessionStatus string
	WSURL         template.JS
	WorkerToken   template.JS
}

// newSessionData builds sessionData with properly typed JS values.
func newSessionData(workerID, workerLabel, workerURL, sessionName, sessionStatus, wsURL, workerToken string) sessionData {
	b, _ := json.Marshal(wsURL)
	tb, _ := json.Marshal(workerToken)
	return sessionData{
		WorkerID:      workerID,
		WorkerLabel:   workerLabel,
		WorkerURL:     workerURL,
		SessionName:   sessionName,
		SessionStatus: sessionStatus,
		WSURL:         template.JS(b),
		WorkerToken:   template.JS(tb),
	}
}
