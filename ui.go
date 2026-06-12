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
  .new-btn{margin-left:.5rem;padding:.5rem 1.1rem;border-radius:8px;
    font-size:.85rem;font-weight:700;font-family:var(--mono);cursor:pointer;
    background:var(--accent-dim);color:var(--accent);
    border:1px solid rgba(88,166,255,.3);transition:background .15s}
  .new-btn:hover{background:rgba(88,166,255,.25)}
  .split-btn{margin-left:auto;padding:.5rem 1rem;border-radius:8px;
    font-size:.85rem;font-weight:700;font-family:var(--mono);cursor:pointer;
    background:rgba(255,255,255,.04);color:var(--text-secondary);text-decoration:none;
    border:1px solid var(--border);transition:all .15s;display:inline-flex;align-items:center;gap:.35rem}
  .split-btn:hover{border-color:var(--border-hi);color:var(--text)}

  .section-label{font-size:.75rem;color:var(--text-secondary);text-transform:uppercase;
    letter-spacing:.1em;margin-bottom:.75rem;font-weight:500}

  @media (max-width:480px){
    /* toolbar: label on its own line, buttons fill the row */
    .toolbar{flex-wrap:wrap;gap:.4rem;margin-bottom:1.25rem}
    .section-label{width:100%;margin-bottom:0}
    .split-btn{margin-left:0;flex:1;justify-content:center;padding:.6rem .5rem}
    .new-btn{margin-left:0;flex:1;padding:.6rem .5rem}
    /* machine tabs: a bit more padding for touch */
    .machine-tab{padding:.3rem .8rem;font-size:.78rem}
  }

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

  .dir-browser{border:1px solid var(--border);border-radius:6px;overflow:hidden;margin-top:.25rem}
  .dir-nav{display:flex;align-items:center;gap:.4rem;padding:.35rem .6rem;
    background:rgba(255,255,255,.04);border-bottom:1px solid var(--border);font-size:.72rem;
    color:var(--text-secondary);min-height:28px}
  .dir-nav-up{background:none;border:none;color:var(--accent);font-family:var(--mono);
    font-size:.72rem;cursor:pointer;padding:0 .3rem;flex-shrink:0}
  .dir-nav-up:disabled{opacity:.3;cursor:default}
  .dir-nav-path{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1;direction:rtl;text-align:left}
  .dir-list{max-height:160px;overflow-y:auto}
  .dir-item{display:flex;align-items:center;gap:.5rem;padding:.38rem .7rem;cursor:pointer;
    font-size:.8rem;transition:background .1s;border:none;background:none;
    width:100%;text-align:left;color:var(--text);font-family:var(--mono)}
  .dir-item:hover{background:rgba(255,255,255,.06)}
  .dir-item.selected{background:var(--accent-dim);color:var(--accent)}
  .dir-item .di-icon{font-size:.75rem;flex-shrink:0;color:var(--text-secondary)}
  .dir-item.git .di-icon{color:var(--green)}
  .dir-item .di-name{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
  .dir-item .di-arrow{font-size:.65rem;color:var(--text-secondary);flex-shrink:0}
  .dir-empty{padding:.6rem .7rem;font-size:.75rem;color:var(--text-secondary)}

  .toast{position:fixed;top:1rem;left:50%;transform:translateX(-50%);
    padding:.55rem 1.1rem;border-radius:8px;font-size:.82rem;font-weight:600;
    font-family:var(--mono);z-index:200;opacity:0;pointer-events:none;
    transition:opacity .2s;max-width:calc(100vw - 2rem);text-align:center}
  .toast.show{opacity:1}
  .toast.error{background:#1c0f0f;border:1px solid rgba(248,81,73,.4);color:var(--red)}
  .toast.info{background:rgba(255,255,255,.06);border:1px solid var(--border-hi);color:var(--text)}
</style>
</head>
<body>

<div id="toast" class="toast"></div>
<div class="title">
  <div class="machine-tabs" id="machine-tabs">
    <button class="machine-tab active" id="tab-all" onclick="setTab('all')">All</button>
  </div>
</div>

<div class="toolbar">
  <div class="section-label">Sessions</div>
  <a class="split-btn" href="/multi">&#9638; Split view</a>
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
        <input id="m-dir" type="text" placeholder="~" oninput="browseDirFromInput(this.value)" autocomplete="off" spellcheck="false">
        <div class="dir-browser" id="m-dirbrowser">
          <div class="dir-nav">
            <button class="dir-nav-up" id="m-dirup" onclick="browseDirUp()" title="Up">↑</button>
            <span class="dir-nav-path" id="m-dirpath">~</span>
          </div>
          <div class="dir-list" id="m-dirlist"><div class="dir-empty">Loading…</div></div>
        </div>
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
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
}

var toastTimer = null;
function showToast(msg, type) {
  var el = document.getElementById('toast');
  el.textContent = msg;
  el.className = 'toast ' + (type || 'error') + ' show';
  clearTimeout(toastTimer);
  toastTimer = setTimeout(function(){ el.classList.remove('show'); }, 3500);
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
    if (!w) {
      activeTab = 'all';
      sessionStorage.setItem('activeTab', 'all');
      document.getElementById('tab-all').classList.add('active');
    } else {
      (w.sessions || []).forEach(function(s){ items.push({s:s, wid:w.id, wlabel:w.label, wurl:w.url}); });
    }
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
    .catch(function(){ document.getElementById('instances').innerHTML = '<div class="empty">Failed to load — retrying…</div>'; });
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
          showToast(d.error || d.message || 'Action failed'); return;
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
  if (!online.length) { showToast('No online workers', 'info'); return; }
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
  browseHistory = [];
  var wid = document.getElementById('m-machine').value;
  browseDir(wid, '~', false);
  document.getElementById('overlay').classList.add('show');
  document.getElementById('m-name').focus();
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

var browseWid = null;
var browseHistory = [];
var browseFetchSeq = 0;
var browseInputTimer = null;

function browseDir(wid, path, fromInput) {
  if (wid) browseWid = wid;
  if (!browseWid) return;
  var dir = (path || '').trim() || '~';
  if (!fromInput) document.getElementById('m-dir').value = dir;
  document.getElementById('m-dirpath').textContent = dir;
  document.getElementById('m-dirlist').innerHTML = '<div class="dir-empty">Loading…</div>';
  document.getElementById('m-dirup').disabled = false;
  var seq = ++browseFetchSeq;
  var url = '/api/workers/' + encodeURIComponent(browseWid) + '/browse?dir=' + encodeURIComponent(dir);
  fetch(url).then(function(r) { return r.json(); }).then(function(d) {
    if (seq !== browseFetchSeq) return; // superseded by a later fetch
    if (d.error) { document.getElementById('m-dirlist').innerHTML = '<div class="dir-empty">' + esc(d.error) + '</div>'; return; }
    var current = d.current || dir;
    if (!fromInput) document.getElementById('m-dir').value = current;
    document.getElementById('m-dirpath').textContent = current;
    document.getElementById('m-dirup').disabled = (current === '/');
    if (fromInput) {
      // Typing navigated to a new location — reset history so Up works correctly.
      browseHistory = [current];
    } else if (!browseHistory.length || browseHistory[browseHistory.length-1] !== current) {
      browseHistory.push(current);
    }
    var dirs = d.dirs || [];
    if (!dirs.length) {
      document.getElementById('m-dirlist').innerHTML = '<div class="dir-empty">No subdirectories</div>';
      return;
    }
    var html = '';
    dirs.forEach(function(item) {
      var cls = 'dir-item' + (item.has_git ? ' git' : '');
      var icon = item.has_git ? '⎇' : '⬡';
      var arrow = item.has_subs ? '›' : '';
      html += '<button class="' + cls + '" data-path="' + esc(item.path) + '" onclick="browseDir(null,this.dataset.path,false)">';
      html += '<span class="di-icon">' + icon + '</span>';
      html += '<span class="di-name">' + esc(item.name) + '</span>';
      html += '<span class="di-arrow">' + arrow + '</span>';
      html += '</button>';
    });
    document.getElementById('m-dirlist').innerHTML = html;
  }).catch(function() {
    if (seq !== browseFetchSeq) return;
    document.getElementById('m-dirlist').innerHTML = '<div class="dir-empty">Failed to load</div>';
  });
}

function browseDirFromInput(val) {
  clearTimeout(browseInputTimer);
  browseInputTimer = setTimeout(function() { browseDir(null, val, true); }, 400);
}

function browseDirUp() {
  if (browseHistory.length > 1) {
    browseHistory.pop();
    var prev = browseHistory[browseHistory.length - 1];
    browseHistory.pop();
    browseDir(null, prev, false);
  } else {
    var cur = document.getElementById('m-dir').value.trim();
    var parent = cur.replace(/\/[^/]+$/, '') || '/';
    browseDir(null, parent, false);
  }
}

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
    if (d.error) { showToast(d.error); return; }
    closeModal();
    var name = d.name || payload.name;
    if (name) location.href = '/session/' + encodeURIComponent(wid) + '/' + encodeURIComponent(name);
    else load();
  })
  .catch(function(){ btn.disabled = false; btn.textContent = 'Spawn'; showToast('Spawn failed'); });
}

document.addEventListener('keydown', function(e){
  if (e.key === 'Escape') closeModal();
  if (e.key === 'Enter' && document.getElementById('overlay').classList.contains('show')) {
    var btn = document.getElementById('m-submit');
    if (!btn.disabled) submitSpawn();
  }
});

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
  border-bottom:1px solid var(--border);z-index:20}
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
  flex-wrap:nowrap;gap:4px;overflow-x:auto;align-items:center;justify-content:center}
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

.toast{position:fixed;top:1rem;left:50%;transform:translateX(-50%);
  padding:.55rem 1.1rem;border-radius:8px;font-size:.82rem;font-weight:600;
  font-family:var(--mono);z-index:200;opacity:0;pointer-events:none;
  transition:opacity .2s;max-width:calc(100vw - 2rem);text-align:center}
.toast.show{opacity:1}
.toast.error{background:#1c0f0f;border:1px solid rgba(248,81,73,.4);color:var(--red)}
.toast.info{background:rgba(13,17,23,.95);border:1px solid rgba(255,255,255,.2);color:var(--text)}
@media (max-width:480px){
  .session-sep,.session-sub{display:none}
  /* bigger tap targets — iOS HIG minimum is 44pt */
  .back{padding:.6rem .7rem .6rem 0;font-size:1.6rem;min-height:44px}
  .menu-btn{padding:10px 14px;min-height:44px;min-width:44px;font-size:1.4rem}
  /* slightly larger session label so it's readable at arm's length */
  .session-label{font-size:.95rem}
  /* compact badge — it's secondary info on mobile */
  .badge{font-size:.6rem;padding:1px 6px;letter-spacing:.02em}
  /* dropdown should be wide enough to tap easily */
  .ditem{padding:.65rem .85rem}
}
</style>
</head>
<body>
<div id="toast" class="toast"></div>
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
    <div class="ak-sep"></div>
    <button class="ak" onclick="pasteToTerm()">paste</button>
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
var statusFetchSeq = 0;  // incremented each call so stale responses are ignored
var reconnectTimer = null;
var reconnectCount = 0;

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

var toastTimer = null;
function showToast(msg, type) {
  var el = document.getElementById('toast');
  el.textContent = msg;
  el.className = 'toast ' + (type || 'error') + ' show';
  clearTimeout(toastTimer);
  toastTimer = setTimeout(function(){ el.classList.remove('show'); }, 3500);
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

function pasteToTerm() {
  if (!navigator.clipboard || !navigator.clipboard.readText) return;
  navigator.clipboard.readText().then(function(text) {
    if (!text) return;
    term.scrollToBottom();
    if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({type:'input', data:text}));
    term.focus();
  }).catch(function(){});
}

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
  ov.style.cssText = 'position:absolute;inset:0;z-index:1;touch-action:none;-webkit-tap-highlight-color:transparent;';
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
    reconnectCount = 0;
    if (WS_TOKEN) ws.send(JSON.stringify({type: 'auth', token: WS_TOKEN}));
    if (INIT_STATUS === 'running') setBadge('live');
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
    reconnectTimer = null;
    // Code 4000 = worker confirmed session is stopped. Show stopped and do not reconnect.
    if (e.code === 4000) { setBadge('stopped'); return; }
    // Stopped sessions: show output once then stop reconnecting.
    if (INIT_STATUS !== 'running') { setBadge('stopped'); return; }
    reconnectCount++;
    if (reconnectCount < 10) {
      setBadge('connecting');
      reconnectTimer = setTimeout(connect, 2000);
    } else {
      setBadge('stopped');
      showToast('Connection lost — refresh to retry', 'info');
    }
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
  // Sequence number ensures only the latest call applies — concurrent calls
  // (page load + ws.onopen) won't overwrite each other with stale data.
  var seq = ++statusFetchSeq;
  fetch('/api/workers/' + encodeURIComponent(WORKER_ID) + '/instances')
    .then(function(r) { return r.json(); })
    .then(function(sessions) {
      if (seq !== statusFetchSeq) return; // superseded by a later call
      var sess = sessions.find(function(s) { return s.name === SESS_NAME; });
      if (!sess) return;
      INIT_STATUS = sess.status;
      updateSessionButtons(sess.status);
      if (sess.status !== 'running') setBadge('stopped');
    })
    .catch(function(){});
}

function sessionAction(action) {
  fetch('/api/workers/' + encodeURIComponent(WORKER_ID) + '/' + action + '/' + encodeURIComponent(SESS_NAME), {method:'POST'})
    .then(function(r) {
      return r.json().then(function(d) {
        if (!r.ok) { showToast(d.error || d.message || action + ' failed'); return; }
        if (action === 'delete') { location.href = '/'; return; }
        if (action === 'kill') {
          INIT_STATUS = 'stopped';
          updateSessionButtons('stopped');
          // Close WS immediately — don't wait for code 4000 from the worker.
          if (ws) { ws.onclose = null; ws.close(); ws = null; }
          setBadge('stopped');
        } else {
          // resume or restart: cancel any pending reconnect timer, then reconnect.
          INIT_STATUS = 'running';
          reconnectCount = 0;
          updateSessionButtons('running');
          setBadge('connecting');
          if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null; }
          if (ws) { ws.onclose = null; ws.close(); ws = null; }
          reconnectTimer = setTimeout(connect, 500);
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

// File upload — routed through dispatch proxy so it works behind Pangolin SSO.
async function uploadFile(file) {
  if (!file) return null;
  var form = new FormData();
  form.append('file', file);
  try {
    var res = await fetch(
      '/api/workers/' + encodeURIComponent(WORKER_ID) + '/upload/' + encodeURIComponent(SESS_NAME),
      {method: 'POST', body: form}
    );
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

// multiTmpl is the split-view page. All session data is fetched client-side;
// the handler serves this shell with no template variables.
var multiTmpl = template.Must(template.New("multi").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="theme-color" content="#0d1117">
<title>split view — dispatch</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/css/xterm.min.css">
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0d1117;--surface:rgba(255,255,255,0.04);--border:rgba(255,255,255,0.1);
  --border-hi:rgba(255,255,255,0.2);--accent:#58a6ff;--accent-dim:rgba(88,166,255,0.15);
  --green:#3fb950;--green-dim:rgba(63,185,80,0.15);--red:#f85149;--red-dim:rgba(248,81,73,0.15);
  --amber:#d29922;--text:#e2e8f0;--text-dim:#94a3b8;--text-muted:#8b949e;
  --mono:'JetBrains Mono','SF Mono',monospace;
}
html,body{height:100%;overflow:hidden;background:var(--bg);color:var(--text);
  font-family:var(--mono);-webkit-font-smoothing:antialiased}

/* layout */
.multi-layout{display:flex;flex-direction:column;height:100dvh}
.multi-topbar{flex-shrink:0;display:flex;align-items:center;gap:.6rem;
  padding:8px 12px;background:rgba(13,17,23,.95);backdrop-filter:blur(20px);
  border-bottom:1px solid var(--border);z-index:20}
.back{color:var(--accent);background:transparent;border:none;font-size:1.4rem;
  cursor:pointer;display:flex;align-items:center;text-decoration:none;
  transition:opacity .15s;flex-shrink:0;line-height:1}
.back:hover{opacity:.7}
.multi-title{font-size:.85rem;font-weight:700;flex:1;color:var(--text-muted)}
.edit-btn{background:rgba(255,255,255,.06);border:1px solid var(--border);color:var(--text-muted);
  font-family:var(--mono);font-size:.78rem;font-weight:600;padding:.32rem .75rem;
  border-radius:6px;cursor:pointer;transition:all .15s;flex-shrink:0}
.edit-btn:hover{border-color:var(--border-hi);color:var(--text)}

/* pane grid */
.pane-grid{flex:1;display:grid;gap:1px;background:var(--border);overflow:hidden;min-height:0}
.pane{display:flex;flex-direction:column;background:var(--bg);min-height:0}

/* pane header */
.pane-header{flex-shrink:0;height:32px;display:flex;align-items:center;gap:.35rem;
  padding:0 .35rem 0 .65rem;background:rgba(13,17,23,.85);
  border-bottom:1px solid var(--border);overflow:visible;position:relative;z-index:10}
.pane-name{font-size:.78rem;font-weight:700;white-space:nowrap;overflow:hidden;
  text-overflow:ellipsis;flex-shrink:1;min-width:0}
.pane-sep{font-size:.78rem;color:var(--text-muted);flex-shrink:0}
.pane-worker{font-size:.78rem;color:var(--text-muted);white-space:nowrap;flex-shrink:0}
.pane-badge{font-size:.6rem;font-weight:700;letter-spacing:.04em;padding:1px 6px;
  border-radius:999px;flex-shrink:0;white-space:nowrap}
.pane-badge.badge-live{background:var(--green-dim);color:var(--green);border:1px solid rgba(52,211,153,.25)}
.pane-badge.badge-stopped{background:rgba(255,255,255,.04);color:var(--text-muted);border:1px solid var(--border)}
.pane-badge.badge-connecting{background:var(--accent-dim);color:var(--accent);border:1px solid rgba(88,166,255,.25)}
.pane-menu-wrap{position:relative;margin-left:auto;flex-shrink:0}
.pane-mbtn{background:transparent;border:none;color:var(--text-dim);cursor:pointer;
  font-size:1.1rem;padding:2px 6px;border-radius:4px;transition:background .12s;line-height:1}
.pane-mbtn:hover{background:var(--surface)}
.pane-dropdown{position:absolute;right:0;top:calc(100% + 2px);background:#0f1018;
  border:1px solid var(--border-hi);border-radius:8px;min-width:145px;overflow:hidden;
  z-index:50;box-shadow:0 8px 24px rgba(0,0,0,.6);display:none}
.pane-dropdown.open{display:block}

/* terminal */
.pane-term{flex:1;overflow:hidden;min-height:0}
.pane-term .xterm{height:100%}
.pane-term .xterm-viewport{overflow-y:auto!important}

/* dropdown items */
.ditem{display:block;width:100%;text-align:left;background:transparent;border:none;
  color:var(--text-dim);font-family:var(--mono);font-size:.8rem;padding:.45rem .8rem;
  cursor:pointer;transition:background .1s,color .1s}
.ditem:hover{background:var(--surface);color:var(--text)}
.ditem.red:hover{color:var(--red)}
.ditem.amber:hover{color:var(--amber)}
.ditem.green:hover{color:var(--green)}
.dsep{height:1px;background:var(--border);margin:2px 0}

/* picker modal */
.overlay{position:fixed;inset:0;background:rgba(0,0,0,.75);backdrop-filter:blur(12px);
  z-index:100;display:flex;align-items:center;justify-content:center;
  opacity:0;pointer-events:none;transition:opacity .2s}
.overlay.show{opacity:1;pointer-events:auto}
.modal{background:#161b22;border:1px solid var(--border-hi);border-radius:12px;
  width:100%;max-width:500px;overflow:hidden;
  transform:translateY(16px);transition:transform .2s cubic-bezier(.16,1,.3,1)}
.overlay.show .modal{transform:translateY(0)}
.modal-head{display:flex;align-items:center;justify-content:space-between;
  padding:.85rem 1.1rem;border-bottom:1px solid var(--border)}
.modal-title{font-size:.9rem;font-weight:700}
.modal-close{background:none;border:none;color:var(--text-muted);font-size:1.2rem;
  cursor:pointer;padding:2px 5px;transition:color .12s}
.modal-close:hover{color:var(--text)}
.modal-body{padding:.75rem 1rem;max-height:56vh;overflow-y:auto;
  display:flex;flex-direction:column;gap:.6rem}
.modal-foot{padding:.7rem 1.1rem;border-top:1px solid var(--border);
  display:flex;align-items:center;gap:.5rem}
.picker-hint{flex:1;font-size:.75rem;color:var(--text-muted)}
.btn{font-size:.82rem;font-weight:700;padding:.42rem .95rem;font-family:var(--mono);
  border-radius:6px;border:none;cursor:pointer;transition:all .12s}
.btn-cancel{background:rgba(255,255,255,.06);color:var(--text-muted);border:1px solid var(--border)}
.btn-cancel:hover{color:var(--text)}
.btn-ok{background:var(--accent-dim);color:var(--accent);border:1px solid rgba(88,166,255,.3)}
.btn-ok:hover:not(:disabled){background:rgba(88,166,255,.25)}
.btn-ok:disabled{opacity:.4;cursor:not-allowed}

/* picker list */
.picker-machine{display:flex;flex-direction:column;gap:2px}
.picker-machine-label{font-size:.7rem;font-weight:700;color:var(--text-muted);
  text-transform:uppercase;letter-spacing:.08em;padding:.3rem .1rem .15rem;
  display:flex;align-items:center;gap:.4rem}
.status-dot{width:6px;height:6px;border-radius:50%;flex-shrink:0;background:var(--border)}
.status-dot.online{background:var(--green)}
.status-dot.offline{background:var(--red)}
.picker-session{display:flex;align-items:center;gap:.6rem;padding:.38rem .6rem;
  border-radius:6px;cursor:pointer;border:1px solid transparent;transition:background .1s}
.picker-session:hover{background:var(--surface)}
.picker-session.checked{background:var(--accent-dim);border-color:rgba(88,166,255,.25)}
.picker-session input[type=checkbox]{accent-color:var(--accent);flex-shrink:0;cursor:pointer}
.picker-sname{flex:1;font-size:.82rem;font-weight:600}
.picker-sbadge{font-size:.65rem;font-weight:600;padding:1px 7px;border-radius:999px;flex-shrink:0}
.picker-sbadge.badge-running{background:var(--green-dim);color:var(--green)}
.picker-sbadge.badge-stopped{background:rgba(255,255,255,.06);color:var(--text-muted)}

/* toast */
.toast{position:fixed;top:1rem;left:50%;transform:translateX(-50%);padding:.55rem 1.1rem;
  border-radius:8px;font-size:.82rem;font-weight:600;font-family:var(--mono);z-index:200;
  opacity:0;pointer-events:none;transition:opacity .2s;max-width:calc(100vw - 2rem);text-align:center}
.toast.show{opacity:1}
.toast.error{background:#1c0f0f;border:1px solid rgba(248,81,73,.4);color:var(--red)}
.toast.info{background:rgba(13,17,23,.95);border:1px solid rgba(255,255,255,.2);color:var(--text)}
</style>
</head>
<body>

<div id="toast" class="toast"></div>

<div class="overlay" id="picker-ov">
  <div class="modal">
    <div class="modal-head">
      <span class="modal-title">Select sessions</span>
      <button class="modal-close" id="picker-close" onclick="closePicker()">&#x2715;</button>
    </div>
    <div class="modal-body" id="picker-body">
      <div style="color:var(--text-muted);font-size:.8rem">Loading&#8230;</div>
    </div>
    <div class="modal-foot">
      <span class="picker-hint" id="picker-hint">Select 2&#8211;4 sessions</span>
      <button class="btn btn-cancel" id="picker-cancel" onclick="closePicker()">Cancel</button>
      <button class="btn btn-ok" id="picker-ok" disabled onclick="submitPicker()">Open</button>
    </div>
  </div>
</div>

<div class="multi-layout">
  <div class="multi-topbar">
    <a class="back" href="/" title="Dashboard">&#8592;</a>
    <span class="multi-title">split view</span>
    <button class="edit-btn" onclick="openPicker()">Edit sessions</button>
  </div>
  <div id="pane-grid" class="pane-grid"></div>
</div>

<script src="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/lib/xterm.js"></script>
<script src="https://cdn.jsdelivr.net/npm/@xterm/addon-fit@0.11.0/lib/addon-fit.js"></script>
<script>
var panesData   = [];  // [{workerID,sessionName,workerLabel,wsURL,token,initStatus}]
var paneStates  = [];  // [{term,fit,ws,reconnectTimer,reconnectCount}]
var pickerSelected = [];
var toastTimer  = null;

// ── bootstrap ─────────────────────────────────────────────────────────────────

(function() {
  var params = new URLSearchParams(location.search);
  var sessions = params.getAll('s');
  if (sessions.length >= 2 && sessions.length <= 4) {
    loadPanesFromURL(sessions);
  } else {
    openPicker();
  }
})();

// ── toast ─────────────────────────────────────────────────────────────────────

function showToast(msg, type) {
  var el = document.getElementById('toast');
  el.textContent = msg;
  el.className = 'toast ' + (type || 'error') + ' show';
  clearTimeout(toastTimer);
  toastTimer = setTimeout(function(){ el.classList.remove('show'); }, 3500);
}

// ── utils ─────────────────────────────────────────────────────────────────────

function esc(s) {
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
}

// ── URL loading ───────────────────────────────────────────────────────────────

function loadPanesFromURL(sessions) {
  fetch('/api/workers')
    .then(function(r){ return r.json(); })
    .then(function(workers) {
      var infos = sessions.map(function(s) {
        var slash = s.indexOf('/');
        var wid = s.substring(0, slash), name = s.substring(slash + 1);
        var w = workers.find(function(wk){ return wk.id === wid; });
        var sess = w ? (w.sessions || []).find(function(ss){ return ss.name === name; }) : null;
        return {workerID:wid, sessionName:name, workerLabel:w ? w.label : wid,
          initStatus:sess ? sess.status : 'stopped', wsURL:'', token:''};
      });
      return Promise.all(infos.map(function(p) {
        return fetch('/api/workers/' + encodeURIComponent(p.workerID) + '/ws/' + encodeURIComponent(p.sessionName))
          .then(function(r){ return r.ok ? r.json() : {}; })
          .then(function(info){ p.wsURL = info.ws_url || ''; p.token = info.token || ''; return p; });
      }));
    })
    .then(renderGrid)
    .catch(function(){ openPicker(); });
}

// ── picker ────────────────────────────────────────────────────────────────────

function openPicker() {
  pickerSelected = panesData.map(function(p) {
    return {key:p.workerID+'/'+p.sessionName, workerID:p.workerID, sessionName:p.sessionName, workerLabel:p.workerLabel};
  });
  var hasPanes = panesData.length > 0;
  document.getElementById('picker-cancel').style.display = hasPanes ? '' : 'none';
  document.getElementById('picker-close').style.display  = hasPanes ? '' : 'none';
  document.getElementById('picker-body').innerHTML = '<div style="color:var(--text-muted);font-size:.8rem">Loading&#8230;</div>';
  document.getElementById('picker-ov').classList.add('show');
  fetch('/api/workers')
    .then(function(r){ return r.json(); })
    .then(renderPickerBody)
    .catch(function(){ document.getElementById('picker-body').innerHTML = '<div style="color:var(--red)">Failed to load sessions</div>'; });
}

function closePicker() {
  if (!panesData.length) return;
  document.getElementById('picker-ov').classList.remove('show');
}

function renderPickerBody(workers) {
  var html = '';
  workers.forEach(function(w) {
    var sessions = w.sessions || [];
    if (!sessions.length) return;
    html += '<div class="picker-machine">';
    html += '<div class="picker-machine-label"><span class="status-dot ' + (w.online ? 'online' : 'offline') + '"></span>' + esc(w.label) + '</div>';
    sessions.forEach(function(s) {
      var key = w.id + '/' + s.name;
      var checked = pickerSelected.some(function(p){ return p.key === key; });
      html += '<label class="picker-session' + (checked ? ' checked' : '') + '" data-key="' + esc(key) + '">';
      html += '<input type="checkbox" class="picker-cb" value="' + esc(key) + '"' + (checked ? ' checked' : '') + ' data-wlabel="' + esc(w.label) + '" onchange="togglePickerSession(this)">';
      html += '<span class="picker-sname">' + esc(s.name) + '</span>';
      html += '<span class="picker-sbadge badge-' + s.status + '">' + s.status + '</span>';
      html += '</label>';
    });
    html += '</div>';
  });
  if (!html) html = '<div style="color:var(--text-muted);font-size:.8rem">No sessions available</div>';
  document.getElementById('picker-body').innerHTML = html;
  updatePickerUI();
}

function togglePickerSession(cb) {
  var key = cb.value;
  var slash = key.indexOf('/');
  var wid = key.substring(0, slash), name = key.substring(slash + 1);
  var label = cb.closest ? cb.closest('.picker-session') : null;
  if (cb.checked) {
    pickerSelected.push({key:key, workerID:wid, sessionName:name, workerLabel:cb.dataset.wlabel || ''});
  } else {
    pickerSelected = pickerSelected.filter(function(p){ return p.key !== key; });
  }
  if (label) label.classList.toggle('checked', cb.checked);
  updatePickerUI();
}

function updatePickerUI() {
  var n = pickerSelected.length;
  var btn = document.getElementById('picker-ok');
  var hint = document.getElementById('picker-hint');
  btn.disabled = n < 2;
  btn.textContent = n >= 2 ? 'Open (' + n + ')' : 'Open';
  hint.textContent = n < 2 ? 'Select 2–4 sessions' : n + ' selected';
  document.querySelectorAll('.picker-cb').forEach(function(cb) {
    if (!cb.checked) cb.disabled = n >= 4;
  });
}

function submitPicker() {
  if (pickerSelected.length < 2) return;
  var selected = pickerSelected.slice();
  // tear down existing panes before closing so closePicker guard works
  tearDownPanes();
  document.getElementById('picker-ov').classList.remove('show');
  document.getElementById('pane-grid').innerHTML = '';

  Promise.all(selected.map(function(p) {
    return fetch('/api/workers/' + encodeURIComponent(p.workerID) + '/ws/' + encodeURIComponent(p.sessionName))
      .then(function(r){ return r.ok ? r.json() : {}; })
      .then(function(info) {
        return {workerID:p.workerID, sessionName:p.sessionName, workerLabel:p.workerLabel,
          wsURL:info.ws_url || '', token:info.token || '', initStatus:'running'};
      });
  })).then(function(paneInfos) {
    return fetch('/api/workers').then(function(r){ return r.json(); }).then(function(workers) {
      paneInfos.forEach(function(p) {
        var w = workers.find(function(wk){ return wk.id === p.workerID; });
        if (w) {
          var sess = (w.sessions || []).find(function(s){ return s.name === p.sessionName; });
          if (sess) p.initStatus = sess.status;
        }
      });
      return paneInfos;
    });
  }).then(renderGrid).catch(function(){ showToast('Failed to open sessions'); openPicker(); });
}

function tearDownPanes() {
  paneStates.forEach(function(st) {
    if (!st) return;
    if (st.ws) { st.ws.onclose = null; st.ws.close(); }
    if (st.reconnectTimer) clearTimeout(st.reconnectTimer);
    if (st.term) st.term.dispose();
  });
  panesData = [];
  paneStates = [];
}

// ── grid ──────────────────────────────────────────────────────────────────────

function renderGrid(paneInfosArr) {
  panesData  = paneInfosArr.slice();
  paneStates = new Array(panesData.length).fill(null);

  var grid = document.getElementById('pane-grid');
  grid.innerHTML = '';
  setGridCSS();

  var params = new URLSearchParams();
  panesData.forEach(function(p){ params.append('s', p.workerID + '/' + p.sessionName); });
  history.replaceState(null, '', '/multi?' + params.toString());

  panesData.forEach(function(p, i){ grid.appendChild(createPaneEl(i, p)); });

  requestAnimationFrame(function(){
    panesData.forEach(function(_, i){ initPane(i); });
  });
}

function setGridCSS() {
  var n = panesData.length;
  var grid = document.getElementById('pane-grid');
  if (n === 2) { grid.style.gridTemplateColumns = '1fr 1fr'; grid.style.gridTemplateRows = '1fr'; }
  else if (n === 3) { grid.style.gridTemplateColumns = '1fr 1fr 1fr'; grid.style.gridTemplateRows = '1fr'; }
  else { grid.style.gridTemplateColumns = '1fr 1fr'; grid.style.gridTemplateRows = '1fr 1fr'; }
}

function createPaneEl(i, p) {
  var running = p.initStatus === 'running';
  var showRun  = running ? '' : 'display:none';
  var showStop = running ? 'display:none' : '';
  var d = document.createElement('div');
  d.className = 'pane'; d.id = 'pane-' + i;
  d.innerHTML =
    '<div class="pane-header">' +
      '<span class="pane-name">' + esc(p.sessionName) + '</span>' +
      '<span class="pane-sep">&#183;</span>' +
      '<span class="pane-worker">' + esc(p.workerLabel) + '</span>' +
      '<span class="pane-badge badge-' + (running ? 'connecting' : 'stopped') + '" id="pbadge-' + i + '">' + (running ? 'connecting' : 'stopped') + '</span>' +
      '<div class="pane-menu-wrap">' +
        '<button class="pane-mbtn" onclick="togglePaneMenu(' + i + ')">&#8943;</button>' +
        '<div class="pane-dropdown" id="pdrop-' + i + '">' +
          '<button class="ditem" id="pd-kill-'    + i + '" style="' + showRun  + '" onclick="paneAction(' + i + ',\'kill\')">Kill</button>' +
          '<button class="ditem amber" id="pd-restart-' + i + '" style="' + showRun  + '" onclick="paneAction(' + i + ',\'restart\')">Restart</button>' +
          '<button class="ditem green" id="pd-resume-'  + i + '" style="' + showStop + '" onclick="paneAction(' + i + ',\'resume\')">Resume</button>' +
          '<button class="ditem red"   id="pd-delete-'  + i + '" style="' + showStop + '" onclick="paneAction(' + i + ',\'delete\')">Delete</button>' +
          '<div class="dsep"></div>' +
          '<button class="ditem" onclick="closePane(' + i + ')">Close pane</button>' +
        '</div>' +
      '</div>' +
    '</div>' +
    '<div class="pane-term" id="pterm-' + i + '"></div>';
  return d;
}

// ── terminal init ─────────────────────────────────────────────────────────────

function initPane(i) {
  var p = panesData[i];
  var termEl = document.getElementById('pterm-' + i);
  if (!p || !termEl) return;

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
    scrollback:2000, convertEol:false,
  });

  var fit = new FitAddon.FitAddon();
  term.loadAddon(fit);
  term.open(termEl);
  fit.fit();

  paneStates[i] = {term:term, fit:fit, ws:null, reconnectTimer:null, reconnectCount:0};

  new ResizeObserver(function(){ fit.fit(); sendPaneResize(i); }).observe(termEl);

  term.onData(function(data) {
    var st = paneStates[i];
    if (st && st.ws && st.ws.readyState === WebSocket.OPEN)
      st.ws.send(JSON.stringify({type:'input', data:data}));
  });

  connectPane(i);
  updatePaneButtons(i, p.initStatus);
}

// ── WS per pane ───────────────────────────────────────────────────────────────

function connectPane(i) {
  var p = panesData[i];
  var st = paneStates[i];
  if (!p || !st || !p.wsURL) return;

  setPaneBadge(i, p.initStatus === 'running' ? 'connecting' : 'stopped');

  var ws = new WebSocket(p.wsURL);
  ws.binaryType = 'arraybuffer';
  st.ws = ws;

  ws.onopen = function() {
    st.reconnectCount = 0;
    if (p.token) ws.send(JSON.stringify({type:'auth', token:p.token}));
    if (p.initStatus === 'running') setPaneBadge(i, 'live');
    sendPaneResize(i);
  };

  ws.onmessage = function(e) {
    if (st.term) st.term.write(new Uint8Array(e.data));
  };

  ws.onclose = function(e) {
    if (paneStates[i] !== st) return; // pane was replaced
    st.ws = null;
    if (e.code === 4000) {
      p.initStatus = 'stopped';
      setPaneBadge(i, 'stopped');
      updatePaneButtons(i, 'stopped');
      return;
    }
    if (p.initStatus !== 'running') { setPaneBadge(i, 'stopped'); return; }
    st.reconnectCount++;
    if (st.reconnectCount < 10) {
      setPaneBadge(i, 'connecting');
      st.reconnectTimer = setTimeout(function(){ connectPane(i); }, 2000);
    } else {
      setPaneBadge(i, 'stopped');
      showToast(p.sessionName + ': connection lost', 'info');
    }
  };

  ws.onerror = function() { setPaneBadge(i, 'connecting'); };
}

function sendPaneResize(i) {
  var st = paneStates[i];
  if (!st || !st.term || !st.ws || st.ws.readyState !== WebSocket.OPEN) return;
  st.ws.send(JSON.stringify({type:'resize', cols:st.term.cols, rows:st.term.rows}));
}

function setPaneBadge(i, state) {
  var el = document.getElementById('pbadge-' + i);
  if (!el) return;
  el.className = 'pane-badge badge-' + state;
  el.textContent = state;
}

function updatePaneButtons(i, status) {
  var running = status === 'running';
  var map = {kill:running, restart:running, resume:!running, delete:!running};
  ['kill','restart','resume','delete'].forEach(function(a) {
    var el = document.getElementById('pd-' + a + '-' + i);
    if (el) el.style.display = map[a] ? '' : 'none';
  });
}

// ── pane actions ──────────────────────────────────────────────────────────────

function paneAction(i, action) {
  var p = panesData[i];
  var st = paneStates[i];
  if (!p) return;
  closePaneMenu(i);
  fetch('/api/workers/' + encodeURIComponent(p.workerID) + '/' + action + '/' + encodeURIComponent(p.sessionName), {method:'POST'})
    .then(function(r){ return r.json().then(function(d) {
      if (!r.ok) { showToast(d.error || action + ' failed'); return; }
      if (action === 'delete') { closePane(i); return; }
      if (action === 'kill') {
        p.initStatus = 'stopped';
        if (st && st.ws) { st.ws.onclose = null; st.ws.close(); st.ws = null; }
        setPaneBadge(i, 'stopped');
        updatePaneButtons(i, 'stopped');
      } else {
        p.initStatus = 'running';
        if (st) st.reconnectCount = 0;
        if (st && st.ws) { st.ws.onclose = null; st.ws.close(); st.ws = null; }
        setPaneBadge(i, 'connecting');
        updatePaneButtons(i, 'running');
        setTimeout(function(){ connectPane(i); }, 500);
      }
    }); })
    .catch(function(){ showToast(action + ' failed'); });
}

function closePane(i) {
  var st = paneStates[i];
  if (st) {
    if (st.ws) { st.ws.onclose = null; st.ws.close(); }
    if (st.reconnectTimer) clearTimeout(st.reconnectTimer);
    if (st.term) st.term.dispose();
  }
  panesData.splice(i, 1);
  paneStates.splice(i, 1);

  if (panesData.length === 1) {
    var p = panesData[0];
    location.href = '/session/' + encodeURIComponent(p.workerID) + '/' + encodeURIComponent(p.sessionName);
    return;
  }
  if (panesData.length === 0) {
    history.replaceState(null, '', '/multi');
    document.getElementById('pane-grid').innerHTML = '';
    openPicker();
    return;
  }

  // Re-render remaining panes (indices shift after splice).
  // Terminals reconnect via replay buffer so no output is lost.
  var remaining = panesData.slice();
  paneStates.forEach(function(s) {
    if (!s) return;
    if (s.ws) { s.ws.onclose = null; s.ws.close(); }
    if (s.reconnectTimer) clearTimeout(s.reconnectTimer);
    if (s.term) s.term.dispose();
  });
  panesData = [];
  paneStates = [];
  document.getElementById('pane-grid').innerHTML = '';
  renderGrid(remaining);
}

// ── pane menu ─────────────────────────────────────────────────────────────────

function togglePaneMenu(i) {
  var dd = document.getElementById('pdrop-' + i);
  if (!dd) return;
  var wasOpen = dd.classList.contains('open');
  closeAllPaneMenus();
  if (!wasOpen) dd.classList.add('open');
}

function closePaneMenu(i) {
  var dd = document.getElementById('pdrop-' + i);
  if (dd) dd.classList.remove('open');
}

function closeAllPaneMenus() {
  document.querySelectorAll('.pane-dropdown.open').forEach(function(d){ d.classList.remove('open'); });
}

document.addEventListener('click', function(e) {
  if (!e.target.closest || !e.target.closest('.pane-menu-wrap')) closeAllPaneMenus();
});

document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    closeAllPaneMenus();
    if (panesData.length > 0) closePicker();
  }
});
</script>
</body>
</html>
`))
