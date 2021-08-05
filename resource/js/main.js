
$(function(){
    var scrollPosition;
    $('.js-modal-open').each(function(){
        $(this).on('click',function(){
            var target = $(this).data('target');
            var modal = document.getElementById(target);
            $(modal).fadeIn();
            scrollPosition = $(window).scrollTop();
            $('body').addClass('fixed').css({'top': -scrollPosition});
            return false;
        });
    });
    $('.js-modal-close').on('click',function(){
        $('.js-modal').fadeOut();
        $('body').removeClass('fixed').css({'top': 0});
		window.scrollTo(0, scrollPosition);
        return false;
    }); 

    $('input[name="visibility"]').change(function() {
         $('.sche-'+$(this).val()).toggle();
    });
});
