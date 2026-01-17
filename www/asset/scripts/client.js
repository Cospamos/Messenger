const Client = (ws) => {
  RunEvents(ws)
  const refs = Refs(); 
  ws.onmessage = (event) => {
    
    try {
      const arr = JSON.parse(event.data)
      if (arr.MessageType == "error") {
        ShowSysInfo(arr.MessageContent[0])
      }
      if (arr.MessageType == "clients") {
        UsersContainer.innerHTML = ''
        arr.MessageContent.forEach((client, idx) => {
          if (client.ID != arr.Addressee) {
            const name = client.Name ? client.Name : client.ID.slice(0, 5); 
            const ico = client.ImgURL ? client.ImgURL : (idx + 1)
            MakeUser(ico, name);
          } else {
            refs.UserName.forEach(el => { // My user name
              const name = client.Name ? client.Name : client.ID.slice(0, 5);
              el.textContent = name;
            })
            if (client.ImgURL) {
              refs.HeadUserIcoImg.forEach(el => { el.src = client.ImgURL })

              refs.HeadUserIcoImg.forEach(el => el.style.display = "block")
              refs.HeadUserIcoText.forEach(el => el.style.display = "none")
            } else {
              refs.HeadUserIcoText.forEach(el => { el.textContent = (idx + 1) })
              refs.HeadUserIcoImg.forEach(el => el.style.display = "none")
              refs.HeadUserIcoText.forEach(el => el.style.display = "block")
            }
          }
        });
      }

      if (arr.MessageType == "messages") {
        MakeMessages(arr.MessageContent, arr.Addressee)
      }
    } catch (err) {
      console.error(err)
    }
   
  }

  ws.onclose = (event) => {
    if (event.code == "1008") { // Policy violation
      alert("Too many requests: ", + event.reason);
    }
  }
}

(async () => {
  const ws = new WebSocket(`http://localhost:8080/ws`);
  Client(ws)
})();

