

function Prompt(){

    let toast = function (c) {

        const {
             //default values
            msg = " ",
            icon = "success",
            postition = "top-end",
        } = c;

        let Toast = Swal.mixin({
            toast: true,
            icon: icon,
            title: msg,
            position: postition,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });
        Toast.fire({});


    }

    let success = function (c) {
        const {
            icon = "success",
            title = "",
            msg = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer,
        });

    }


    let error = function (c) {
        const { 
            
            icon = "error",
            title = "Oops..",
            msg = "Something went wrong!",
            footer = "",
        } = c;

        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer,
        });

    }

    async function custom(c) {
        const { 
            icon ="",
            msg = "",
            title = "",
            showConfirmButton = true,

        } = c;
        const { value: formValues} = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if(c.willOpen !== undefined){
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById("start").value,
                    document.getElementById("end").value
                ];
            },
            didOpen: () => {

                if(c.didOpen !== undefined){
                    c.didOpen();
                }

            }


        });
        if (formValues) { 
           
            if (formValues.dismiss !== Swal.DismissReason.cancel) {
                
                    if (formValues.value !== "") { 
                     
                     if (c.callback !== undefined) { 
                        c.callback(formValues)
                    } else {
                        c.callback(false);
                    }
                } else { 
                    c.callback(false);
                }
            }
        }

    }

    return {
            toast: toast,
            success: success,
            error: error,
            custom: custom,
        }
}


