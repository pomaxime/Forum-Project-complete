document.addEventListener('DOMContentLoaded', function () {
  var params = new URLSearchParams(window.location.search);
  var category = params.get('category') || '';
  var filter = params.get('filter') || '';

  document.querySelectorAll('.filters a').forEach(function (link) {
    var p = new URLSearchParams(link.getAttribute('href').split('?')[1] || '');
    if (p.get('category') === category && p.get('filter') === filter) {
      link.style.fontWeight = 'bold';
      link.style.borderColor = '#0077cc';
      link.style.color = '#0077cc';
    }
  });
});