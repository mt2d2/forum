$(function() {
  // inspired by http://stackoverflow.com/questions/22636819/confirm-delete-using-bootstrap-3-modal-box
  $('button[name="removePost"]').on('click', function(e) {
    var $form = $(this).closest('form');
    e.preventDefault();
    $('#confirm-delete').modal({ keyboard: false })
      .one('click', '#delete', function() {
        $form.trigger('submit');
      });
  });
});
