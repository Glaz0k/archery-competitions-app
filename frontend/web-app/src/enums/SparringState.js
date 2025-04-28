import createEnum from "../helper/createEnum";

const SparringState = createEnum({
  ONGOING: { value: "ongoing", textValue: "" },
  TOP_WIN: { value: "top_win", textValue: "" },
  BOT_WIN: { value: "bot_win", textValue: "" },
});

export default SparringState;
