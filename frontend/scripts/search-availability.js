import { Datepicker } from 'vanillajs-datepicker';
import Swal from 'sweetalert2/dist/sweetalert2.min.js';

const MAX_HOURS = 4; // 使用者最多可以預約幾小時

async function postAvailabilityCheck(formData) {
  console.log("📝 Posting availability check with:", Object.fromEntries(formData.entries()));
  const response = await fetch('/check-availability', {
    method: "POST",
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'X-CSRF-Token': document.querySelector('input[name="csrf_token"]').value
    },
    body: new URLSearchParams(formData),
  });
  
  if (!response.ok) {
    throw new Error("Failed to fetch availability");
  }

  return await response.json();
}

async function postReservation(formData) {
  console.log("📝 make reservation with:", Object.fromEntries(formData.entries()));
  const response = await fetch('/make-reservation', {
    method: "POST",
    credentials: 'same-origin',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'X-CSRF-Token': document.querySelector('input[name="csrf_token"]').value,
    },
    body: new URLSearchParams(formData),
  });

  if (!response.ok) {
    throw new Error("Reservation failed");
  }

  return await response.json();
}

async function loadAvailabilityForDate(selectedDate) {
  console.log("📅 Date selected (from loadAvailabilityForDate):", selectedDate);

  if (!selectedDate) return;

  try {
    const response = await fetch(`/check-availability?date=${selectedDate}`);
    const data = await response.json(); //  [{ hour: 8, available_courts: 4 }, ...]
    console.log("✅ Hourly availability:", data);

    // 清空選單
    startHourSelect.innerHTML = '<option value="" disabled selected hidden>選擇開始時間</option>';
    endHourSelect.innerHTML = '<option value="" disabled selected hidden>選擇結束時間</option>';
    endHourSelect.disabled = true;

    // 將每小時可用數量加進 start/end hour 選單
    data.forEach(slot => {
      const startOpt = document.createElement("option");
      startOpt.value = slot.value;
      startOpt.textContent = slot.label;

      const endOpt = document.createElement("option");
      endOpt.value = slot.value + 1;
      endOpt.textContent = `${slot.value + 1}:00`;

      startHourSelect.appendChild(startOpt);
      endHourSelect.appendChild(endOpt);
    });

  } catch (err) {
    console.error("❌ Failed to fetch availability:", err);
  }
}


// 📅 DOM references
const searchForm = document.getElementById("availability-search");
const endHourSelect = searchForm.elements["end_hour"];
const startHourSelect = searchForm.elements["start_hour"];
const dateInput = searchForm.elements["date"];

// ⏳ 顯示查詢資訊
const availabilityInfo = document.createElement("p");
availabilityInfo.className = "has-text-success mt-3";
searchForm.appendChild(availabilityInfo);

// 📅 日期選擇器
const today = new Date();
const threeMonthsLater = new Date();
threeMonthsLater.setMonth(today.getMonth() + 3);

// 防止超過 12 月時錯亂（例如從 10 月加到 13 月）
if (threeMonthsLater.getMonth() < today.getMonth()) {
  threeMonthsLater.setFullYear(today.getFullYear() + 1);
}
const datepicker = new Datepicker(dateInput, {
  format: "yyyy-mm-dd",
  autohide: true,
  minDate: today, //禁用今天之前的日期
  maxDate: threeMonthsLater, // 限制日期選擇在三個月後
});
datepicker.setDate(datepicker.getFocusedDate());

// 🕒 選擇開始時間後，動態載入結束時間
startHourSelect.addEventListener("change", async () => {
  const date = dateInput.value;
  const startHourRaw = startHourSelect.value;

  console.log("📍 startHour 選擇:", startHourRaw);

  // 加入原始值檢查，避免 undefined 被送出
  if (!date || !startHourRaw || startHourRaw === "undefined") {
    console.warn("⚠️ 缺少日期或開始時間不正確");
    endHourSelect.disabled = true;
    endHourSelect.innerHTML = '<option value="" disabled selected hidden>選擇結束時間</option>';
    return;
  }

  try {
    const res = await fetch(`/check-availability?date=${date}&start_hour=${startHourRaw}`);
    if (!res.ok) throw new Error("❌ fetch 失敗");

    const data = await res.json();
    console.log("📊 hour-availability 資料:", data);

    // 篩掉 <= startHour 的時段
    const startHour = parseInt(startHourSelect.value, 10); // 加上 parseInt
    const filtered = data.filter(d => d.value > startHour && d.value <= startHour + MAX_HOURS);
    console.log("🔍 可選結束時段 (hour > start):", filtered);

    // 清空 endHour 選單
    endHourSelect.innerHTML = '<option value="" disabled selected hidden>選擇結束時間</option>';

    console.table(data);  // 印出整張 hour availability 表格
    console.log("選擇的開始時間:", startHour);

    if (filtered.length === 0) {
      const option = document.createElement("option");
      option.textContent = "⚠️ 沒有可選結束時間";
      option.disabled = true;
      endHourSelect.appendChild(option);
      endHourSelect.disabled = true;
    }
    
    // 動態填入選項
    filtered.forEach(d => {
      const option = document.createElement("option");
      option.value = d.value;
      option.textContent = d.label;
      endHourSelect.appendChild(option);
    });
    endHourSelect.disabled = false;

  } catch (err) {
    console.error("❌ 查詢或處理錯誤:", err);
    endHourSelect.disabled = true;
    endHourSelect.innerHTML = '<option value="" disabled selected hidden>載入失敗</option>';
  }
});

// 📩 送出預約表單
searchForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const formData = new FormData(searchForm);
  const dateStr = formData.get("date");
  const startRawLabel = startHourSelect.selectedOptions[0].text;
  const endRawLabel = endHourSelect.selectedOptions[0].text;

  // 移除多餘文字（可用 xx 面）
  const startLabel = startRawLabel.replace(/（可用.*?）/, "").trim();
  const endLabel = endRawLabel.replace(/（可用.*?）/, "").trim();

  // 星期幾
  const weekdays = ["日", "一", "二", "三", "四", "五", "六"];
  const dateObj = new Date(dateStr);
  const weekday = weekdays[dateObj.getDay()];

  // 固定場地數與金額（你可日後換成變數）
  const courtCount = 1;
  const price = 400;

  // SweetAlert 確認視窗
  const result = await Swal.fire({
    title: "確認預約",
    html: `📅 日期：${dateStr}（星期${weekday}）<br>🕒 時間：${startLabel} - ${endLabel}<br>🏸 場地數：${courtCount} 面<br>💰 金額：${price} 元`,
    showCancelButton: true,
    confirmButtonText: "預約",
    cancelButtonText: "取消",
  });

  if (!result.isConfirmed) return;

  try {
    const res = await postReservation(formData);

    // 成功視窗也顯示金額與場地
    Swal.fire({
      icon: "success",
      title: "預約成功！",
      html: `📅 日期：${res.date}<br>🕒 ${res.start_hour} - ${res.end_hour} 點<br>🏸 場地數：${courtCount} 面<br>💰 金額：${price} 元`,
    });
  } catch (err) {
    Swal.fire({
      icon: "error",
      text: "預約失敗，請稍後再試",
    });
  }
});

