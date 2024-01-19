// moved back into base.alout.html

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
    const { title = "Def Mult Inputs", html = "" } = c;

    const { value: result } = await Swal.fire({
      title: title,
      html: html,
      focusConfirm: true,
      backdrop: true,
      showCancelButton: true,
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
