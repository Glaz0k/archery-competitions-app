import createEnum from "../helper/createEnum";

const BowClass = createEnum({
  CLASSIC: { value: "classic", textValue: "Классический лук" },
  BLOCK: { value: "block", textValue: "Блочный лук" },
  CLASSIC_NEWBIE: { value: "classic_newbie", textValue: "Классический лук (новички)" },
  CLASSIC_3D: { value: "3D_classic", textValue: "3Д-классический лук" },
  COMPOUND_3D: { value: "3D_compound", textValue: "3Д-составной лук" },
  LONG_3D: { value: "3D_long", textValue: "3Д-длинный лук" },
  PERIPHERAL: { value: "peripheral", textValue: "Периферийный лук" },
  PERIPHERAL_WITH_RING: {
    value: "peripheral_with_ring",
    textValue: "Периферийный лук (с кольцом)",
  },
});

export default BowClass;
