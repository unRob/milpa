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

commandSelector.addEventListener("keyup", function(evt){
  let cmd = this.value
  if (evt.keyCode == 13 || evt.keyCode == 9) {
    if (cmd != initialValue && allCommands.includes(cmd)) {
      window.location = `/${cmd.replaceAll(" ", "/")}/`
    } else {
      return true
    }
  }
})

