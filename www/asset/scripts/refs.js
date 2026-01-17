const $ = (selector) => document.querySelector(selector);
const $$ = (selector) => document.querySelectorAll(selector);

Element.prototype.on = function (event, handler) {
  this.addEventListener(event, handler);
  return this;
};

NodeList.prototype.on = function (event, handler) {
  this.forEach(el => el.addEventListener(event, handler));
  return this;
};

Window.prototype.on = function (event, handler) {
  this.addEventListener(event, handler);
  return this;
};

Document.prototype.on = function (event, handler) {
  this.addEventListener(event, handler);
  return this;
};

const Refs = () => ({
  // Buttons
  SendButton: $(".send_btn"),
  ActionsBtn: $$(".users_actions_btn"),
  UserSetingsBtn: $$(".user_settings_btn"),
  ShowQRbtns: $$(".show_qr"),
  SaveUsersData: $(".save_users_data"),
  // Containers
  UserModalContainer: $$(".user.modal_container"),
  ChangeLightMode: $$(".ligh_mode"),
  GlobalModalContainer: $(".global.user_container"),
  QRContainer: $(".qr_container"),
  UserContainer: $(".user_settings_container"),
  Messages: $$(".message"),
  UserName: $$(".head_user_name"),
  // Icons
  HeadUserIcon: $(".head_user_ico"),
  UserIconOverview: $(".overview img"),
  HeadUserIcoText: $$(".head_user_ico p"),
  HeadUserIcoImg: $$(".head_user_ico img"),
  // Inputs
  MessageInput: $(".message_input"),
  ChangeUserIconInput: $(".change_user_ico input"),
  ChangeUserNameInput: $(".change_user_name input"),
});
