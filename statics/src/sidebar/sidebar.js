import './sidebar.css';
/* global bootstrap: false */

(function () {
  'use strict'
  var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
  tooltipTriggerList.forEach(function (tooltipTriggerEl) {
    new bootstrap.Tooltip(tooltipTriggerEl)
  })
})()

let SidebarSwitch = function (){
  const location = window.location.href
  const id = location.substring(location.lastIndexOf("/"))
  document.getElementById(id).classList.add("active")
}

window.addEventListener('load', SidebarSwitch);
