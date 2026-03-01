function toArr(s){ return s.split(',').map(x=>x.trim()).filter(x=>x); }

async function post(path, body){
  const res = await fetch(path, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(body)});
  return res.json();
}

document.getElementById('btnEval').addEventListener('click', async ()=>{
  const hole = toArr(document.getElementById('evalHole').value);
  const community = toArr(document.getElementById('evalCommunity').value);
  const out = document.getElementById('evalOut'); out.textContent = 'Loading...';
  try{ const r = await post('/evaluate', {hole, community}); out.textContent = JSON.stringify(r, null, 2); }catch(e){ out.textContent = 'Error: '+e }
});

document.getElementById('btnCompare').addEventListener('click', async ()=>{
  const h1 = { hole: toArr(document.getElementById('c1h').value), community: toArr(document.getElementById('c1c').value) };
  const h2 = { hole: toArr(document.getElementById('c2h').value), community: toArr(document.getElementById('c2c').value) };
  const out = document.getElementById('compareOut'); out.textContent = 'Loading...';
  try{ const r = await post('/compare', {hand1: h1, hand2: h2}); out.textContent = JSON.stringify(r, null, 2); }catch(e){ out.textContent = 'Error: '+e }
});

document.getElementById('btnSim').addEventListener('click', async ()=>{
  const hero = { hole: toArr(document.getElementById('simHero').value) };
  const community = toArr(document.getElementById('simComm').value);
  const num = parseInt(document.getElementById('simPlayers').value)||2;
  const iter = parseInt(document.getElementById('simIter').value)||1000;
  const out = document.getElementById('simOut'); out.textContent = 'Running...';
  try{ const r = await post('/simulate', { hero, community_known: community, num_players: num, iterations: iter, concurrency: 4 }); out.textContent = JSON.stringify(r, null, 2); }catch(e){ out.textContent = 'Error: '+e }
});
