function toggleComments(id) {
  var el = document.getElementById('comments-' + id);
  if (el) el.style.display = el.style.display === 'none' ? 'block' : 'none';
}

function copyShareLink(button) {
  var relativeUrl = button.dataset.shareUrl;
  if (!relativeUrl) {
    return;
  }

  var fullUrl = window.location.origin + relativeUrl;

  if (navigator.share) {
    navigator.share({
      title: document.title,
      url: fullUrl,
    }).catch(function () {
      copyToClipboard(fullUrl);
    });
    return;
  }

  copyToClipboard(fullUrl);
}

function copyToClipboard(text) {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).then(function () {
      alert('Lien copié dans le presse-papiers');
    }).catch(function () {
      window.prompt('Copiez le lien du post', text);
    });
    return;
  }

  window.prompt('Copiez le lien du post', text);
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