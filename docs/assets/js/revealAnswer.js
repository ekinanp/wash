function revealAnswer(exercise) {
  var answer = document.getElementById(exercise).getElementsByClassName("answer")[0];
  // Calculating the display via answer.style.display would make it so that
  // users would have to double click the button before revealing the answer.
  // See https://stackoverflow.com/questions/21852932/javascript-onclick-requires-two-clicks
  var display = window.getComputedStyle(answer ,null).getPropertyValue("display")
  if (display === "none") {
    answer.style.display = "block";
  } else {
    answer.style.display = "none";
  }
}
