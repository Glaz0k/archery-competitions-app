import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';
import 'widgets/user.dart';

class EditProfilePage extends StatefulWidget {
  const EditProfilePage({super.key});

  @override
  State<EditProfilePage> createState() => _EditProfilePage();
}

class _EditProfilePage extends State<EditProfilePage> {
  late TextEditingController fullNameController;
  late TextEditingController regionController;
  late TextEditingController clubController;
  String? chosenGender;
  String? chosenBow;
  String? chosenRank;
  late DateTime chosenTime;

  @override
  void initState() {
    super.initState();
    final user = Provider.of<UserProvider>(context, listen: false).userPref;
    fullNameController = TextEditingController(text: user.fullName);
    regionController = TextEditingController(text: user.region);
    clubController = TextEditingController(text: user.club);
    chosenGender = user.identity.getGender;
    chosenBow = user.bow?.getBowClass;
    chosenRank = user.rank?.getSportsRank;
    chosenTime = user.dateOfBirth;
  }

  @override
  void dispose() {
    fullNameController.dispose();
    regionController.dispose();
    clubController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final userProvider = Provider.of<UserProvider>(context);
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: OnionBar("Редактирование профиля", context),
          body: KeyboardDismisser(
            gestures: [GestureType.onTap],
            child: GestureDetector(
              child: SingleChildScrollView(
                child: Container(
                  padding: const EdgeInsets.all(8),
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: Column(
                          children: [
                            buildMyField(),
                            SizedBox(height: 7),
                            DropdownButtonFormField<String>(
                              value: chosenGender,
                              decoration: InputDecoration(
                                labelText: 'Пол',
                                border: OutlineInputBorder(),
                              ),
                              items:
                                  ['Мужчина', 'Женщина'].map((String value) {
                                    return DropdownMenuItem<String>(
                                      value: value,
                                      child: Text(value),
                                    );
                                  }).toList(),
                              onChanged: (newVal) {
                                setState(() {
                                  chosenGender = newVal;
                                });
                              },
                            ),
                            SizedBox(height: 7),
                            DropdownButtonFormField<String>(
                              value: chosenBow,
                              decoration: InputDecoration(
                                labelText: 'Класс лука',
                                border: OutlineInputBorder(),
                              ),
                              items:
                                  [
                                    'Классический',
                                    'Блочный',
                                    "КЛ(новички)",
                                    "3Д-классический лук",
                                    "3Д-составной лук",
                                    "3Д-длинный лук",
                                    "Периферийный лук",
                                    "Периферийный лук(с кольцом)",
                                  ].map((String value) {
                                    return DropdownMenuItem<String>(
                                      value: value,
                                      child: Text(value),
                                    );
                                  }).toList(),
                              onChanged: (newVal) {
                                setState(() {
                                  chosenBow = newVal;
                                });
                              },
                            ),
                            SizedBox(height: 7),
                            DropdownButtonFormField<String>(
                              value: chosenRank,
                              decoration: InputDecoration(
                                labelText: 'Спортивный разряд/звание',
                                border: OutlineInputBorder(),
                              ),
                              items:
                                  [
                                    'Мастер спорта',
                                    'Кандидат мастера спорта',
                                    "Первый разряд",
                                    "Второй разряд",
                                    "Третий разряд",
                                    "Международный магистр",
                                    "Периферийный лук",
                                    "Заслуженный мастер спорта",
                                  ].map((String value) {
                                    return DropdownMenuItem<String>(
                                      value: value,
                                      child: Text(value),
                                    );
                                  }).toList(),
                              onChanged: (newVal) {
                                setState(() {
                                  chosenRank = newVal;
                                });
                              },
                            ),
                            SizedBox(height: 7),

                            SizedBox(height: 7),
                            buildRegion(),
                            SizedBox(height: 7),
                            buildClub(),
                            SizedBox(height: 7),
                            SizedBox(
                              width: double.infinity,
                              // Растягивает на всю ширину
                              child: OutlinedButton(
                                onPressed: () => chooseDate(context),
                                child: Text(
                                  'Изменить дату рождения',
                                  style: TextStyle(
                                    fontSize: 13,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ),
                            ),
                            SizedBox(height: 7),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                              children: [
                                ElevatedButton(
                                  onPressed: () {
                                    Navigator.pop(context);
                                  },
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: Colors.redAccent,
                                    side: BorderSide.none,
                                    shape: StadiumBorder(),
                                  ),
                                  child: const Text(
                                    "Отменить",
                                    style: TextStyle(color: Colors.white),
                                  ),
                                ),
                                ElevatedButton(
                                  onPressed: () {
                                    userProvider.updateFullName(
                                      fullNameController.text,
                                    );
                                    userProvider.updateRegion(
                                      regionController.text,
                                    );
                                    userProvider.updateClub(
                                      clubController.text,
                                    );
                                    chosenGender == "Мужчина"
                                        ? userProvider.updateGender(Gender.male)
                                        : userProvider.updateGender(
                                          Gender.female,
                                        );
                                    userProvider.updateBow(
                                      BowExtension.setBowClass(chosenBow),
                                    );
                                    userProvider.updateRank(
                                      SportsRankExtension.setSportsRank(
                                        chosenRank,
                                      ),
                                    );
                                    userProvider.updateDateOfBirth(chosenTime);
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        showCloseIcon: true,
                                        behavior: SnackBarBehavior.floating,
                                        backgroundColor: Colors.green,
                                        content: Container(
                                          decoration: BoxDecoration(
                                            color: Colors.green,
                                            shape: BoxShape.circle,
                                          ),
                                          child: Row(
                                            children: [
                                              Icon(
                                                Icons.check_circle,
                                                // color: Colors.white,
                                                size: 40,
                                              ),
                                              SizedBox(width: 30),
                                              Expanded(
                                                child: Text(
                                                  "Изменения сохранены!",
                                                  style: TextStyle(
                                                    fontSize: 20,
                                                    fontWeight: FontWeight.bold,
                                                  ),
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                      ),
                                    );
                                    Navigator.pop(context);
                                  },
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: Colors.green,
                                    side: BorderSide.none,
                                    shape: StadiumBorder(),
                                  ),
                                  child: const Text(
                                    "Сохранить",
                                    style: TextStyle(color: Colors.white),
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future chooseDate(BuildContext context) async {
    final newDate = await showDatePicker(
      context: context,
      initialDate: chosenTime,
      firstDate: DateTime(1950),
      lastDate: DateTime.now(),
    );
    if (newDate == null) return;
    setState(() {
      chosenTime = newDate;
    });
  }

  Widget buildMyField() {
    return TextFormField(
      controller: fullNameController,
      decoration: InputDecoration(
        label: Text("Фамилия Имя", style: TextStyle()),
        prefixIcon: Icon(Icons.person),
        border: OutlineInputBorder(),
      ),
      textInputAction: TextInputAction.done,
      keyboardType: TextInputType.text,
    );
  }

  Widget buildRegion() {
    return TextFormField(
      controller: regionController,
      decoration: InputDecoration(
        label: Text("Регион", style: TextStyle()),
        prefixIcon: Icon(Icons.apartment),
        border: OutlineInputBorder(),
      ),
      textInputAction: TextInputAction.done,
      keyboardType: TextInputType.text,
    );
  }

  Widget buildClub() {
    return TextFormField(
      controller: clubController,
      decoration: InputDecoration(
        label: Text("Клуб", style: TextStyle()),
        prefixIcon: Icon(Icons.school),
        border: OutlineInputBorder(),
      ),
      textInputAction: TextInputAction.done,
      keyboardType: TextInputType.text,
    );
  }
}
