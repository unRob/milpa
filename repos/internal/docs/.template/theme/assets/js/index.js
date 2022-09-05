(function navigateOnCommand () {
  const allCommands = Array.from(document.querySelectorAll("#milpa-commands option")).map((opt) => opt.value)
  const commandSelector = document.querySelector("#command-selector")
  const initialValue = commandSelector.value

  commandSelector.addEventListener("keydown", function(evt){
    let cmd = this.value
    if (evt.keyCode == 9 && cmd != initialValue && allCommands.includes(cmd)) {
      evt.preventDefault()
      return
    }
  })

  function navigateOnChange(input) {
    let cmd = input.value
    if (cmd != initialValue && allCommands.includes(cmd)) {
      window.location = `/${cmd.replaceAll(" ", "/")}/`
    } else {
      return true
    }
  }

  commandSelector.addEventListener("keyup", function(evt){
    if (evt.keyCode == 13 || evt.keyCode == 9) {
     return navigateOnChange(this)
    }
  })

  commandSelector.addEventListener('change', function(evt) {
    return navigateOnChange(this)
  })

  const sidebar = document.querySelector("#sidebar");
  document.querySelector("#menu-toggle").addEventListener("click", function(evt){
    let classes = this.classList
    classes.toggle("open")
    if (classes.contains("open")) {
      sidebar.classList.add("open")
      sidebar.setAttribute("aria-hidden", "false")
      document.body.classList.add("sidebar-open")
    } else {
      sidebar.classList.remove("open")
      sidebar.setAttribute("aria-hidden", "true")
      document.body.classList.remove("sidebar-open")
    }
  })

  document.querySelectorAll('#table-of-contents a').forEach(link => link.addEventListener('click', function(evt) {
    evt.preventDefault()
    document.querySelector("#menu-toggle").classList.remove('open')
    sidebar.classList.remove("open")
    sidebar.setAttribute("aria-hidden", "true")
    document.body.classList.remove("sidebar-open")
    document.getElementById(this.hash.substr(1)).parentNode.scrollIntoView({behavior: "smooth", block: "start"})
    history.pushState(null, null, this.hash)
    return false
  }))
})();
