// Makes navbar burger toggles menu
const $navbarBurger = document.querySelector('.navbar-burger');
$navbarBurger.addEventListener('click', () => {
  // Get the target from the "data-target" attribute
  const target = $navbarBurger.dataset.target;
  const $target = document.getElementById(target);
  // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
  $navbarBurger.classList.toggle('is-active');
  $target.classList.toggle('is-active');
});
// TODO fix: toggle behavior of dropdown in mobile screen size
// document.querySelectorAll('.navbar-item.has-dropdown').forEach($item => {
//   $item.addEventListener('click', () => $item.classList.toggle('is-active'))
// })