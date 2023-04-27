$(document).ready(function(){
    inputPath = document.querySelector('#inputPath');
    myfile = document.querySelector('#myfile');
    editor = document.querySelector('#editor');
    btnMenuMain = document.querySelector('#btn-menu-main')
    btnMenuLogin = document.querySelector('#btn-menu-login')
    btnMenuReport = document.querySelector('#btn-menu-reports')


    function readFile(file, callback){
        const reader = new FileReader();
        reader.addEventListener('load', (event) => {
            const content = event.target.result;
            callback(content)
        });
          
        reader.readAsText(file);
    }

    myfile.addEventListener('change', (e) => {
        inputPath.value = e.target.files[0].name;
        readFile(e.target.files[0], (content) => {
            editor.value = content;
        })
    });

    btnMenuMain.addEventListener('click', (e) => {
        $('#main').siblings().slideUp(500);
        $('#main').slideDown(500);
    });

    btnMenuLogin.addEventListener('click', (e) => {
        $('#login-page').siblings().slideUp(500);
        $('#login-page').slideDown(500);
    });

    btnMenuReport.addEventListener('click', (e) => {
        $('#report-page').siblings().slideUp(500);
        $('#report-page').slideDown(500);
    });

});