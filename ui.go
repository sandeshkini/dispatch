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
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="theme-color" content="#0d1117">
<meta name="apple-mobile-web-app-capable" content="yes">
<title>dispatch</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0d1117;
  --surface:rgba(255,255,255,0.04);
  --border:rgba(255,255,255,0.1);
  --border-hi:rgba(255,255,255,0.2);
  --accent:#58a6ff;
  --accent-dim:rgba(88,166,255,0.15);
  --green:#3fb950;
  --green-dim:rgba(63,185,80,0.15);
  --red:#f85149;
  --red-dim:rgba(248,81,73,0.15);
  --text:#e2e8f0;
  --text-secondary:#8b949e;
  --mono:'JetBrains Mono','SF Mono','Consolas',monospace;
}
html,body{min-height:100%;background:var(--bg);color:var(--text);
  font-family:var(--mono);font-size:14px;line-height:1.5;-webkit-font-smoothing:antialiased}

.wrap{max-width:900px;margin:0 auto;padding:1.25rem 1rem 4rem}

.title{display:flex;align-items:center;justify-content:space-between;
  font-size:1rem;font-weight:700;margin-bottom:2rem;
  padding-bottom:1rem;border-bottom:1px solid var(--border)}

/* worker section */
.worker{margin-bottom:2rem}
.worker-hdr{display:flex;align-items:center;justify-content:space-between;
  margin-bottom:.6rem}
.worker-meta{display:flex;align-items:center;gap:.5rem}
.wlabel{font-size:.75rem;font-weight:600;color:var(--text-secondary);
  text-transform:uppercase;letter-spacing:.07em}

.dot{width:5px;height:5px;border-radius:50%;flex-shrink:0}
.dot.live{background:var(--green);animation:pdot 2s ease-in-out infinite}
.dot.off{background:var(--text-secondary)}
.dot.run{background:var(--green)}
.dot.stop{background:rgba(255,255,255,0.15)}
@keyframes pdot{0%,100%{opacity:1}50%{opacity:.3}}

.btn-spawn{font-size:.75rem;font-weight:700;padding:3px 9px;border-radius:6px;border:none;
  background:var(--accent-dim);color:var(--accent);cursor:pointer;
  font-family:var(--mono);transition:background .12s;white-space:nowrap}
.btn-spawn:hover:not(:disabled){background:rgba(88,166,255,.25)}
.btn-spawn:disabled{opacity:.3;cursor:not-allowed}

/* session rows */
.sess-row{display:flex;align-items:center;gap:1rem;
  padding:.7rem .9rem;border-radius:8px;
  background:var(--surface);border:1px solid var(--border);
  margin-bottom:.4rem;transition:border-color .12s;cursor:default;
  animation:cardin .35s cubic-bezier(.16,1,.3,1) both}
.sess-row:hover{border-color:var(--border-hi)}
@keyframes cardin{from{opacity:0;transform:translateY(6px)}to{opacity:1;transform:translateY(0)}}

.sess-left{display:flex;align-items:center;gap:.6rem;flex-shrink:0;width:140px}
.sess-name{font-size:.9rem;font-weight:700;
  white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.sess-summary{flex:1;font-size:.78rem;color:var(--text-secondary);
  white-space:nowrap;overflow:hidden;text-overflow:ellipsis;min-width:0}

.sess-btns{display:flex;gap:4px;opacity:0;transition:opacity .12s;flex-shrink:0}
.sess-row:hover .sess-btns{opacity:1}
.ibtn{width:26px;height:26px;border-radius:6px;border:1px solid var(--border);
  background:none;color:var(--text-secondary);display:flex;align-items:center;
  justify-content:center;cursor:pointer;transition:all .12s}
.ibtn:hover{background:var(--surface);border-color:var(--border-hi);color:var(--text)}
.ibtn.danger:hover{color:var(--red);border-color:rgba(248,81,73,.3);background:var(--red-dim)}
.ibtn.open-btn:hover{color:var(--accent);border-color:rgba(88,166,255,.3);background:var(--accent-dim)}
.ibtn svg{width:13px;height:13px;stroke:currentColor;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}

.no-sess{font-size:.8rem;color:var(--text-secondary);padding:.5rem 0}

.empty{display:flex;flex-direction:column;align-items:center;gap:.75rem;
  padding:4rem 2rem;text-align:center;
  border:1px dashed var(--border);border-radius:8px;margin:2rem 0}
.empty-lbl{font-size:.85rem;color:var(--text-secondary);line-height:1.7}
.empty-lbl code{font-size:.8rem;background:var(--surface);
  padding:1px 6px;border-radius:4px;color:var(--text)}

.overlay{position:fixed;inset:0;background:rgba(0,0,0,.7);
  backdrop-filter:blur(10px);-webkit-backdrop-filter:blur(10px);
  z-index:100;display:flex;align-items:center;justify-content:center;
  opacity:0;pointer-events:none;transition:opacity .2s}
.overlay.show{opacity:1;pointer-events:auto}
.modal{background:#161b22;border:1px solid var(--border-hi);border-radius:8px;
  width:100%;max-width:400px;overflow:hidden;
  transform:translateY(16px);transition:transform .2s cubic-bezier(.16,1,.3,1)}
.overlay.show .modal{transform:translateY(0)}
.modal-head{display:flex;align-items:center;justify-content:space-between;
  padding:.85rem 1.1rem;border-bottom:1px solid var(--border)}
.modal-title{font-size:.9rem;font-weight:700}
.modal-close{background:none;border:none;color:var(--text-secondary);
  font-size:1.2rem;cursor:pointer;line-height:1;padding:2px 4px;transition:color .12s}
.modal-close:hover{color:var(--text)}
.modal-body{padding:1.1rem;display:flex;flex-direction:column;gap:.85rem}
.fg{display:flex;flex-direction:column;gap:.35rem}
.fg label{font-size:.7rem;font-weight:600;color:var(--text-secondary);
  text-transform:uppercase;letter-spacing:.06em}
.fg select,.fg input{background:rgba(255,255,255,.06);border:1px solid var(--border);
  color:var(--text);font-family:var(--mono);font-size:.85rem;
  padding:.45rem .65rem;border-radius:6px;outline:none;
  transition:border-color .15s;width:100%}
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
<div class="wrap">

  <div class="title">
    <span class="title-label">dispatch</span>
  </div>

  <div id="main"></div>

</div>

<!-- spawn modal -->
<div class="overlay" id="overlay" onclick="overlayClick(event)">
  <div class="modal">
    <div class="modal-head">
      <span class="modal-title" id="modal-title">New session</span>
      <button class="modal-close" onclick="closeModal()">&#x2715;</button>
    </div>
    <div class="modal-body">
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
var _wid = '';

function esc(s) {
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function render(workers) {
  var el = document.getElementById('main');
  if (!workers.length) {
    el.innerHTML = '<div class="empty">' +
      '<div class="empty-lbl">No workers registered yet.<br>' +
      'Set <code>hub_url</code> in a worker\'s config to register.</div></div>';
    return;
  }
  var html = '';
  for (var i = 0; i < workers.length; i++) {
    var w = workers[i];
    var sessions = w.sessions || [];
    var caps = w.capabilities || [];

    html += '<div class="worker">';
    html += '<div class="worker-hdr">';
    html += '<div class="worker-meta">';
    html += '<span class="dot ' + (w.online ? 'live' : 'off') + '"></span>';
    html += '<span class="wlabel">' + esc(w.label) + '</span>';
    html += '</div>';
    var dis = w.online ? '' : ' disabled';
    var capStr = caps.join(',');
    html += '<button class="btn-spawn"' + dis + ' onclick="openModal(\'' + w.id + '\',\'' +
      esc(w.label).replace(/'/g,"\\'") + '\',\'' + capStr + '\')">+ New session</button>';
    html += '</div>';

    if (sessions.length) {
      for (var j = 0; j < sessions.length; j++) {
        var s = sessions[j];
        var sdelay = (j * 0.04) + 's';
        html += '<div class="sess-row" style="animation-delay:' + sdelay + '">';
        html += '<div class="sess-left">';
        html += '<span class="dot ' + (s.status === 'running' ? 'run' : 'stop') + '"></span>';
        html += '<span class="sess-name">' + esc(s.name) + '</span>';
        html += '</div>';
        html += '<span class="sess-summary">' + esc(s.summary || s.dir || '') + '</span>';
        html += '<div class="sess-btns">';
        html += '<button class="ibtn open-btn" title="Open" onclick="openSession(\'' + w.id + '\',\'' + esc(s.name) + '\')">';
        html += '<svg viewBox="0 0 24 24"><polyline points="15 3 21 3 21 9"/><path d="M10 14L21 3"/><path d="M21 3H9"/><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5"/></svg>';
        html += '</button>';
        html += '<button class="ibtn danger" title="Kill" onclick="killSession(\'' + w.id + '\',\'' + esc(s.name) + '\')">';
        html += '<svg viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>';
        html += '</button>';
        html += '</div>';
        html += '</div>';
      }
    } else {
      html += '<div class="no-sess">No sessions</div>';
    }

    html += '</div>';
  }
  el.innerHTML = html;
}

function load() {
  fetch('/api/workers')
    .then(function(r){ return r.json(); })
    .then(function(d){ render(d || []); })
    .catch(function(){ });
}

function openSession(wid, name) {
  location.href = '/session/' + encodeURIComponent(wid) + '/' + encodeURIComponent(name);
}

function killSession(wid, name) {
  if (!confirm('Kill session "' + name + '"?')) return;
  fetch('/api/workers/' + encodeURIComponent(wid) + '/kill/' + encodeURIComponent(name), { method:'POST' })
    .then(function(){ load(); })
    .catch(function(){ alert('Kill failed'); });
}

function openModal(wid, label, capStr) {
  _wid = wid;
  document.getElementById('modal-title').textContent = 'New session on ' + label;
  var sel = document.getElementById('m-cli');
  sel.innerHTML = '<option value="terminal">terminal</option>';
  if (capStr) {
    var caps = capStr.split(',');
    for (var i = 0; i < caps.length; i++) {
      var c = caps[i].trim();
      if (c && c !== 'terminal') {
        sel.innerHTML += '<option value="' + esc(c) + '">' + esc(c) + '</option>';
      }
    }
  }
  document.getElementById('m-dir').value = '';
  document.getElementById('m-name').value = '';
  document.getElementById('overlay').classList.add('show');
  document.getElementById('m-dir').focus();
}

function closeModal() {
  document.getElementById('overlay').classList.remove('show');
}

function overlayClick(e) {
  if (e.target === document.getElementById('overlay')) closeModal();
}

function submitSpawn() {
  var btn = document.getElementById('m-submit');
  btn.disabled = true;
  btn.textContent = 'Spawning...';
  var payload = {
    cli: document.getElementById('m-cli').value,
    dir: document.getElementById('m-dir').value || '~',
    name: document.getElementById('m-name').value || ''
  };
  fetch('/api/workers/' + encodeURIComponent(_wid) + '/spawn', {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(payload)
  })
  .then(function(r){ return r.json(); })
  .then(function(d){
    btn.disabled = false;
    btn.textContent = 'Spawn';
    if (d.error) { alert(d.error); return; }
    var name = d.name || payload.name;
    closeModal();
    if (name) location.href = '/session/' + encodeURIComponent(_wid) + '/' + encodeURIComponent(name);
    else load();
  })
  .catch(function(){
    btn.disabled = false;
    btn.textContent = 'Spawn';
    alert('Spawn failed');
  });
}

document.addEventListener('keydown', function(e){
  if (e.key === 'Escape') closeModal();
});

load();
setInterval(load, 30000);

fetch('/health').then(function(r){ return r.json(); }).then(function(){ }).catch(function(){});
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
  --bg:#07080d;--surface:rgba(255,255,255,0.04);--border:rgba(255,255,255,0.08);
  --accent:#a78bfa;--green:#34d399;--green-dim:rgba(52,211,153,0.12);
  --red:#f87171;--red-dim:rgba(248,113,113,0.1);--amber:#fbbf24;
  --text:#e2e8f0;--text-dim:#94a3b8;--text-muted:#475569;
  --mono:'JetBrains Mono','SF Mono',monospace;
  --sat:env(safe-area-inset-top,0px);--sab:env(safe-area-inset-bottom,0px);
  --sal:env(safe-area-inset-left,0px);--sar:env(safe-area-inset-right,0px);
}
html,body{height:100%;overflow:hidden;background:var(--bg);color:var(--text);
  font-family:var(--mono);-webkit-font-smoothing:antialiased;-webkit-text-size-adjust:100%}
.layout{display:flex;flex-direction:column;height:100%;height:100dvh}

/* topbar */
.topbar{flex-shrink:0;display:flex;align-items:center;gap:.4rem;
  padding:calc(var(--sat) + 8px) calc(var(--sar) + 12px) 8px calc(var(--sal) + 6px);
  background:rgba(7,8,13,.9);backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);
  border-bottom:1px solid var(--border);z-index:10}
.back{color:var(--accent);background:transparent;border:none;
  font-size:1.4rem;cursor:pointer;padding:2px 6px 2px 2px;
  display:flex;align-items:center;line-height:1;text-decoration:none;
  transition:opacity .15s;flex-shrink:0}
.back:hover{opacity:.7}
.session-info{flex:1;min-width:0;display:flex;flex-direction:column;align-items:center;gap:1px}
.session-label{font-size:.85rem;font-weight:700;
  white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:200px}
.session-sub{font-size:.68rem;color:var(--text-muted)}
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
  background:rgba(7,8,13,.9);backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);
  border-top:1px solid var(--border);
  flex-wrap:nowrap;gap:4px;overflow-x:auto;align-items:center}
.abar::-webkit-scrollbar{display:none}
.abar.visible{display:flex}
.ak{height:30px;min-width:40px;padding:0 8px;border-radius:6px;
  background:var(--surface);border:1px solid var(--border);
  color:var(--text-dim);font-family:var(--mono);font-size:.72rem;font-weight:500;
  cursor:pointer;display:flex;align-items:center;justify-content:center;
  transition:background .1s;-webkit-tap-highlight-color:transparent;touch-action:manipulation;flex-shrink:0}
.ak:active{background:rgba(255,255,255,.1)}
.ak.enter{background:rgba(167,139,250,.1);border-color:rgba(167,139,250,.3);color:var(--accent)}
.ak.ctrlc{background:var(--red-dim);border-color:rgba(248,113,113,.25);color:var(--red)}
.ak-sep{width:1px;height:16px;background:var(--border);flex-shrink:0;margin:0 2px}
</style>
</head>
<body>
<div class="layout">

  <div class="topbar">
    <a class="back" href="/" title="Back">&#8592;</a>
    <div class="session-info">
      <div class="session-label">{{.SessionName}}</div>
      <div class="session-sub">{{.WorkerLabel}}</div>
    </div>
    <span class="badge badge-connecting" id="badge">connecting</span>
    <div class="menu-wrap">
      <button class="menu-btn" id="menu-btn" onclick="toggleMenu()">&#8943;</button>
      <div class="dropdown" id="dropdown">
        {{if eq .SessionStatus "running"}}
        <button class="ditem" onclick="sessionAction('kill');closeMenu()">Kill session</button>
        <button class="ditem" onclick="sessionAction('restart');closeMenu()">Restart</button>
        {{else}}
        <button class="ditem" onclick="sessionAction('resume');closeMenu()">Resume</button>
        {{end}}
        <div class="dsep"></div>
        <button class="ditem" onclick="sendCtrlC();closeMenu()">Interrupt (^C)</button>
      </div>
    </div>
  </div>

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
  };
  ws.onmessage = function(e) {
    var data = new Uint8Array(e.data);
    var buf = term.buffer.active;
    var atBot = buf.viewportY + term.rows >= buf.length - 3;
    term.write(data, function(){ if(atBot) term.scrollToBottom(); });
  };
  ws.onclose = function() {
    setBadge(INIT_STATUS === 'running' ? 'connecting' : 'stopped');
    ws = null;
    setTimeout(connect, 2000);
  };
  ws.onerror = function() { setBadge('connecting'); };
}

connect();

function sessionAction(action) {
  fetch('/api/workers/' + encodeURIComponent(WORKER_ID) + '/' + action + '/' + encodeURIComponent(SESS_NAME), {method:'POST'})
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
</script>
</body>
</html>
`))

// sessionData is passed to sessionTmpl.
type sessionData struct {
	WorkerID      string
	WorkerLabel   string
	SessionName   string
	SessionStatus string
	WSURL         template.JS
	WorkerToken   template.JS
}

// newSessionData builds sessionData with properly typed JS values.
func newSessionData(workerID, workerLabel, sessionName, sessionStatus, wsURL, workerToken string) sessionData {
	b, _ := json.Marshal(wsURL)
	tb, _ := json.Marshal(workerToken)
	return sessionData{
		WorkerID:      workerID,
		WorkerLabel:   workerLabel,
		SessionName:   sessionName,
		SessionStatus: sessionStatus,
		WSURL:         template.JS(b),
		WorkerToken:   template.JS(tb),
	}
}
