// function getDataUri(url, callback) {
//
//   let xhr = new XMLHttpRequest();
//   xhr.open("GET", url);
//   xhr.responseType = "arraybuffer";
//   xhr.onload = function()
//   {
//       let buffer = xhr.response;
//       let binary = '';
//       let bytes = new Uint8Array(buffer);
//
//       for (let i = 0, l = bytes.byteLength; i < l; i++) {
//        binary += String.fromCharCode(bytes[i]);
//       }
//       console.log(url, this.getResponseHeader('content-type'));
//       callback("data:" + this.getResponseHeader('content-type') + ";base64," + window.btoa(binary));
//
//   }
//   xhr.send();
//
// }
//
// function iterator(currentValue, index, array) {
//
//   let url = currentValue.getAttribute("src");
//
//   getDataUri(url, function(dataUri) {
//     currentValue.setAttribute("src", dataUri);
//   })
//
// }
//
// $$(".attachment img").forEach(iterator);
//$$(".attachment video source").forEach(iterator);


console.log(JSON.stringify([].concat(...$$(".conversation").map(function(currentValue, index, array) {

  if (currentValue.getAttribute("class").indexOf("placeholder") > -1) {
    return [];
  }

  let group = currentValue.querySelector(".conversation-name").innerHTML;

  return Array.prototype.slice.call(currentValue.querySelectorAll(".message-list li")).map(function(currentValue, index, array) {

    let id = currentValue.getAttribute("id");
    let sender = currentValue.querySelector(".sender").innerHTML;
    let content = (currentValue.querySelector(".content .body") || (document.createElement("span"))).innerHTML;
    let timestamp = currentValue.querySelector(".timestamp").getAttribute("data-timestamp");
    let sent = currentValue.getAttribute("class").indexOf("incoming") === -1;

    let attachments = [];
    currentValue.querySelectorAll(".attachment img").forEach(function(imgElt, index, array) {
      attachments.push({
        kind: "img",
        url: imgElt.getAttribute("src")
      });
    });
    currentValue.querySelectorAll(".attachment video source").forEach(function(videoElt, index, array) {
      attachments.push({
        kind: "video",
        url: videoElt.getAttribute("src")
      });
    });

    if (currentValue.querySelectorAll(".attachment").length != attachments.length) {
      console.error("Unmached attachement", currentValue)
      throw "Unmached attachement"
    }

    return {
      id: id,
      sender: sender,
      content: content,
      timestamp: timestamp,
      sent: sent,
      attachments: attachments,
      kind: "signal",
      group: group
    };

  });

}))));
