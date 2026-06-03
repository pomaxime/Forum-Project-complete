function toggleComments(id) {
  var el = document.getElementById('comments-' + id);
  if (el) el.style.display = el.style.display === 'none' ? 'block' : 'none';
}

document.addEventListener('DOMContentLoaded', function () {
  ['username', 'email', 'password'].forEach(function (fieldName) {
    var field = document.querySelector('input[name="' + fieldName + '"]');
    if (!field) return;
    field.addEventListener('input', function () {
      this.classList.remove('input-error');
    });
  });
});