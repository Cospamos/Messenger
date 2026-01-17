const MessageContainer = document.querySelector(".chat_container .content");
const UsersContainer = document.querySelector(".users_container .online_info")

const MakeMessages = (messagesArr, addressee) => {
  const fragment = document.createDocumentFragment();

  messagesArr.forEach(el => {
    const assigned = (el.Id == addressee) ? "mine" : "other";

    const MessageContent = document.createElement("div");
    MessageContent.classList.add("message_content", assigned);

    if (el.ImgUrl) {
      const UserIco = document.createElement("img");
      UserIco.className = "chat_user_ico";
      UserIco.src = el.ImgUrl;
      MessageContent.appendChild(UserIco);
    } else {
      const UserIco = document.createElement("div");
      UserIco.className = "chat_user_ico";
      UserIco.textContent = assigned[0].toUpperCase();
      MessageContent.appendChild(UserIco);
    }

    const Message = document.createElement("div");
    Message.className = "message";
    Message.textContent = el.Content;
    MessageContent.appendChild(Message);

    fragment.appendChild(MessageContent);
  });

  MessageContainer.appendChild(fragment);
  MessageContainer.scrollTop = MessageContainer.scrollHeight;
};


const MakeUser = (user_ico, user_name) => {
  const UserContent = document.createElement("div");
  UserContent.className = "user_content";
  UsersContainer.append(UserContent)


  if (!isValidImageUrl(user_ico)) { // function from events.js
    const UserIco = document.createElement("div");
    UserIco.className = "user_ico";
    UserIco.textContent = user_ico;
    UserContent.append(UserIco);
  } else {
    const UserIco = document.createElement("img");
    UserIco.className = "user_ico";
    UserIco.src = user_ico;
    UserContent.append(UserIco);
  }

  const UserName = document.createElement("div");
  UserName.className = "user_name";
  UserName.textContent = user_name;
  UserContent.append(UserName);
}

ShowSysInfo = (content) => {
  $(".sys_info_container").classList.add("show")
  $(".sys_info_text_content").textContent = content

  setTimeout(() => {
    $(".sys_info_container").classList.remove("show")
  }, 3000)
}
