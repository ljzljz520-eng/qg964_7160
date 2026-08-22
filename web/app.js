const summary = document.getElementById('summary');
const todo = document.getElementById('todo');
async function loadTodo() {
  const response = await fetch('/todo', { headers: { 'X-Actor': 'demo-supervisor' } });
  if (!response.ok) { summary.textContent = 'Unable to load scoped todo'; return; }
  const records = await response.json();
  summary.textContent = `${records.length} remediation records`;
  todo.replaceChildren(...records.map((record) => {
    const item = document.createElement('li');
    item.textContent = `${record.store_id}: ${record.title} (${record.severity})`;
    return item;
  }));
}
loadTodo();
