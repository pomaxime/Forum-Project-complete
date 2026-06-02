function toggleComments(id) {
  var el = document.getElementById('comments-' + id);
  if (el) el.style.display = el.style.display === 'none' ? 'block' : 'none';
}