function toggleComments(id) {
  var el = document.getElementById('comments-' + id);
  if (el) el.style.display = el.style.display === 'none' ? 'block' : 'none';
}

document.addEventListener('DOMContentLoaded', function () {
  ['username', 'email', 'password', 'title', 'content', 'category'].forEach(function (fieldName) {
    var field = document.querySelector('[name="' + fieldName + '"]');
    if (!field) return;
    field.addEventListener('input', function () {
      this.classList.remove('input-error');
    });
    field.addEventListener('change', function () {
      this.classList.remove('input-error');
    });
  });
});