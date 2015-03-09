$(function() {
  function fadeOut(element) {
    element.animate({ margin: 'hide', padding: 'hide', height: 'hide' }, 'fast', 'swing', function() {
        $(this).remove()
    });
  }

  // inspired by http://stackoverflow.com/questions/22636819/confirm-delete-using-bootstrap-3-modal-box
  $('button[name="removePost"]').on('click', function(e) {
    var $form = $(this).closest('form');
    e.preventDefault();
    $('#confirm-delete').modal({ keyboard: false })
      .one('click', '#delete', function() {
        $form.trigger('submit');
      });
  });

  var alertSuccess = $(".alert-success");
  window.setTimeout(function() {
    fadeOut(alertSuccess)
  }, 2000);
  $('button[class="close"]').on('click', function(e) {
    fadeOut(alertSuccess);
  });
});
