const RunEvents = (ws) => {
  let MessageValue;
  const refs = Refs()
  
  refs.MessageInput.on('keydown', e => {
    if (e.key === 'Enter') {
      e.preventDefault();
      MessageValue = refs.MessageInput.textContent.trim();
      if (MessageValue) {
        let data = {type: "POST", content: MessageValue}
        ws.send(JSON.stringify(data));
        refs.MessageInput.textContent = "";
      }
    }
  });

  refs.SendButton.on("click", () => {
    const value = refs.MessageInput.textContent.trim();
    if (value) {
      const data = { type: "POST", content: value }
      ws.send(JSON.stringify(data));
      refs.MessageInput.textContent = "";
    }
  });

  refs.UserIconOverview.onerror = () => {
    refs.UserIconOverview.style.display = "none";
    refs.UserIconOverview.src = "";
  }

  refs.HeadUserIcoImg.forEach(el => {
    el.onerror = () => {
      el.style.display = "none";
      el.src = "";
    }
  })

  refs.ChangeUserIconInput.on("input", () => {
    const url = refs.ChangeUserIconInput.value.trim();

    if (isValidImageUrl(url)) {
      refs.UserIconOverview.style.display = "block";
      refs.UserIconOverview.src = url;
    } else {
      refs.UserIconOverview.style.display = "none";
      refs.UserIconOverview.src = "";
    }
  });

  refs.ChangeUserIconInput.on("keydown", (e) => {
    if (e.key == "Enter") {
      e.preventDefault()
      const value = refs.ChangeUserIconInput.value.trim();
      
      if (value) {
        const data = {type: "PATCH/ico", content: value}
        ws.send(JSON.stringify(data))
        refs.GlobalModalContainer.classList.add("hidden");
      }
    }
  })

  refs.ChangeUserNameInput.on("keydown", (e) => {
    if (e.key != "Enter") return;
    const value = refs.ChangeUserNameInput.value.trim();  
    if (value) {
      const data = { type: "PATCH/name", content: value }
      ws.send(JSON.stringify(data))
      refs.GlobalModalContainer.classList.add("hidden");
    }
  })

  refs.SaveUsersData.on("click", () => {
    const NameValue = refs.ChangeUserNameInput.value.trim();
    const IconValue = refs.ChangeUserIconInput.value.trim();

    if (NameValue) {
      const data = { type: "PATCH/name", content: NameValue }
      ws.send(JSON.stringify(data))
      refs.GlobalModalContainer.classList.add("hidden");
    }
    if (isValidImageUrl(IconValue)) {
      const data = { type: "PATCH/ico", content: IconValue }
      ws.send(JSON.stringify(data))
      refs.GlobalModalContainer.classList.add("hidden");
    }
  })

  // ActionBtn logic for users container
  refs.ActionsBtn[0].on("click", (e) => {  
    e.stopPropagation()
    refs.UserModalContainer[0].classList.toggle("hidden");
    setTimeout(() => {
      document.addEventListener("click", handleDocClick)
    }, 0)
  })

  // ActionBtn logic for chat container
  refs.ActionsBtn[1].on("click", (e) => {
    e.stopPropagation()
    refs.UserModalContainer[1].classList.toggle("hidden");
    setTimeout(() => {
      document.addEventListener("click", handleDocClick)
    }, 0)
  })

  function handleDocClick(e) {
    if (!refs.UserModalContainer[0].contains(e.target)) {
      refs.UserModalContainer[0].classList.add("hidden")
      document.removeEventListener("click", handleDocClick)
    }

    if (!refs.UserModalContainer[1].contains(e.target)) {
      refs.UserModalContainer[1].classList.add("hidden")
      document.removeEventListener("click", handleDocClick)
    }
  }

  refs.ShowQRbtns.on("click", () => {
    refs.GlobalModalContainer.classList.remove("hidden")
    refs.QRContainer.classList.remove("hidden")
  })
  
  refs.GlobalModalContainer.on("click", (e) => {
    if (e.target !== e.currentTarget) return; 
    refs.GlobalModalContainer.classList.add("hidden");
    refs.QRContainer.classList.add("hidden");
    refs.UserContainer.classList.add("hidden");
  });

  refs.UserSetingsBtn.on("click", () => {
    refs.GlobalModalContainer.classList.remove("hidden")
    refs.UserContainer.classList.remove("hidden");
  })

  refs.ChangeLightMode.on("click", () => {
    document.documentElement.classList.toggle("light");
    refs.ActionsBtn[0].classList.toggle("light"); // ActionBtn from users container
    refs.ActionsBtn[1].classList.toggle("light"); // ActionBtn from chat container
    refs.MessageInput.classList.toggle("light")
    refs.SendButton.classList.toggle("light");
    refs.HeadUserIcon.classList.toggle("light");
    refs.Messages.forEach(el => el.classList.toggle("light"))
  })
}
