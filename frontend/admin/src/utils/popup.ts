import Swal, { SweetAlertIcon } from 'sweetalert2'

export const successPopup = (
  title: string,
  text: string = '',
  confirmButtonText: string = 'Tamam'
) => {
  return Swal.fire({
    title,
    text,
    icon: 'success',
    confirmButtonText
  });
}

export const errorPopup = (
  title: string,
  text: string = '',
  confirmButtonText: string = 'Tamam'
) => {
  return Swal.fire({
    title,
    text,
    icon: 'error',
    confirmButtonText
  });
};

export const infoPopup = (
  title: string,
  text: string = '',
  icon: SweetAlertIcon = 'info',
  confirmButtonText: string = 'Tamam'
) => {
  return Swal.fire({
    title,
    text,
    icon,
    confirmButtonText
  });
};

export const confirmPopup = async (
  title: string,
  text: string,
  confirmButtonText: string = 'Evet',
  cancelButtonText: string = 'Hayır',
  icon: SweetAlertIcon = 'warning'
): Promise<boolean> => {
  const result = await Swal.fire({
    title,
    text,
    icon,
    showCancelButton: true,
    confirmButtonColor: "#3085d6",
    cancelButtonColor: "#d33",
    confirmButtonText,
    cancelButtonText
  });

  return result.isConfirmed;
};

export const customButtonsPopup = async (
  title: string,
  text: string,
  confirmText: string = 'Evet, sil!',
  cancelText: string = 'Vazgeç',
  icon: SweetAlertIcon = 'warning'
): Promise<'confirmed' | 'cancelled'> => {
  const swalWithBootstrapButtons = Swal.mixin({
    customClass: {
      confirmButton: "btn btn-success",
      cancelButton: "btn btn-danger"
    },
    buttonsStyling: false
  });

  const result = await swalWithBootstrapButtons.fire({
    title,
    text,
    icon,
    showCancelButton: true,
    confirmButtonText: confirmText,
    cancelButtonText: cancelText,
    reverseButtons: true
  });

  if (result.isConfirmed) {
    return 'confirmed';
  } else {
    return 'cancelled';
  }
};

export const inputPopup = async (
  title: string,
  inputLabel: string = '',
  placeholder: string = '',
  confirmButtonText: string = 'Gönder',
  cancelButtonText: string = 'İptal',
  inputType: 'text' | 'email' | 'number' | 'password' = 'text'
): Promise<string | null> => {
  const result = await Swal.fire({
    title,
    input: inputType,
    inputLabel,
    inputPlaceholder: placeholder,
    showCancelButton: true,
    confirmButtonText,
    cancelButtonText
  });

  return result.isConfirmed ? result.value : null;
};

