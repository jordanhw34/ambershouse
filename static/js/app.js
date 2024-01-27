// Example starter JavaScript for disabling form submissions if there are invalid fields
// Leaving this on base.layout.html so is available to any/all forms
(() => {
  "use strict";
  // Fetch all the forms we want to apply custom Bootstrap validation styles to
  const forms = document.querySelectorAll(".needs-validation");
  // Loop over them and prevent submission
  Array.from(forms).forEach((form) => {
    form.addEventListener(
      "submit",
      (event) => {
        if (!form.checkValidity()) {
          event.preventDefault();
          event.stopPropagation();
        }

        form.classList.add("was-validated");
      },
      false
    );
  });
})();

function Prompt() {
  let toast = (c) => {
    const {
      title = "",
      icon = "success",
      position = "top-end", // top-start, top-center, top-end
    } = c;

    const Toast = Swal.mixin({
      toast: true,
      title: title,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.onmouseenter = Swal.stopTimer;
        toast.onmouseleave = Swal.resumeTimer;
      },
    });
    Toast.fire({});
  };

  let success = (c) => {
    const { title = "", text = "" } = c;
    Swal.fire({
      icon: "success",
      title: title,
      text: text,
    });
  };

  let error = (c) => {
    const { title = "", text = "" } = c;
    Swal.fire({
      icon: "error",
      title: title,
      text: text,
    });
  };

  let custom = async (c) => {
    const {
      title = "Def Mult Inputs",
      html = "",
      showConfirmButton = true,
      showCancelButton = true,
    } = c;

    const { value: result } = await Swal.fire({
      title: title,
      html: html,
      focusConfirm: true,
      backdrop: true,
      showConfirmButton: showConfirmButton,
      showCancelButton: showCancelButton,
      preConfirm: () => {
        return [
          document.getElementById("start").value,
          document.getElementById("end").value,
        ];
      },
      willOpen: () => {
        if (c.willOpen !== undefined) {
          c.willOpen();
        }
      },
      didOpen: () => {
        if (c.didOpen !== undefined) {
          c.didOpen();
        }
      },
    });

    if (result) {
      //Swal.fire(JSON.stringify(formValues));
      if (result.dismiss === Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (c.callback !== undefined) {
            c.callback(result);
          }
        } else {
          c.callback(false);
        }
      } else {
        c.callback(false);
      }
    }
  };

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  };
}

function notify(msgType, msg) {
  notie.alert({
    //type: Number|String, // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    type: msgType,
    text: msg,
    stay: false,
    time: 2,
    position: "top",
  });
}

function notifyModal(title, text, icon, confirmButtonText) {
  Swal.fire({
    title: title,
    text: text,
    icon: icon,
    confirmButtonText: confirmButtonText,
  });
}

