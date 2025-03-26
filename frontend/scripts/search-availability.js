import { Datepicker } from 'vanillajs-datepicker';
import Swal from 'sweetalert2/dist/sweetalert2.min.js'

async function postReservation(formData) {
  const response = await fetch('/search-availability', {
    method: "post",
    body: new URLSearchParams(formData),
  })
  if (response.status !== 201) {
    let message = undefined;
    try { message = (await response.json()).message } catch (e) { }
    throw new Error(message, {
      cause: `POST request failed (status ${response.status})`
    })
  }
  return await response.json()
}

const searchForm = document.getElementById("availability-search")

// TODO disable non-reservable dates
const datepicker = new Datepicker(
  searchForm.elements.date,
  { format: "yyyy-mm-dd", autohide: true }
);
datepicker.setDate(datepicker.getFocusedDate())

searchForm.addEventListener("submit", async (e) => {
  // prevents <form> action
  e.preventDefault()

  // show confirmation
  const result = await Swal.fire({
    title: "確認預約",
    text: `已選擇 日期：${searchForm.elements.date.value}，時段：${searchForm.elements.timeslot.selectedOptions[0].text}`,
    showCancelButton: true,
    confirmButtonText: "預約",
    cancelButtonText: "取消",
  })
  if (!result.isConfirmed)
    return

  // post reservation and show result
  try {
    const data = await postReservation(new FormData(searchForm))
    Swal.fire({
      icon: "success",
      title: "預約成功！",
      html: `📅 日期：${data.date}<br>🕒 時間：${data.time}`
    });
  } catch (error) {
    Swal.fire({
      icon: "error",
      text: "發生錯誤，請稍後再試！"
    });
    console.log(error);
  }
});