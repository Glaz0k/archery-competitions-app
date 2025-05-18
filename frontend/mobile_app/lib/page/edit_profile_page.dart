import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import '../api/common.dart';
import 'main_competition_page.dart';
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
  late TextEditingController federationController;
  Gender? chosenGender;
  BowClass? chosenBow;
  SportsRank? chosenRank;
  String birthDate = '';

  @override
  void initState() {
    super.initState();
    final api = Provider.of<Api>(context, listen: false);
    final user =
        Provider.of<UserProvider>(context, listen: false).getUser(api)!;
    fullNameController = TextEditingController(text: user.fullName);
    regionController = TextEditingController(text: user.region);
    clubController = TextEditingController(text: user.club);
    federationController = TextEditingController(text: user.federation);
    chosenGender = user.identity;
    chosenBow = user.bow;
    chosenRank = user.rank;
    birthDate = user.birthDate;
  }

  @override
  void dispose() {
    fullNameController.dispose();
    regionController.dispose();
    clubController.dispose();
    federationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final api = Provider.of<Api>(context);
    final userProvider = context.watch<UserProvider>();
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: OnionBar.withoutProfile("Редактирование профиля", context),
          body: SingleChildScrollView(
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
                        DropdownButtonFormField<Gender>(
                          value: chosenGender,
                          decoration: const InputDecoration(
                            labelText: 'Пол',
                            border: OutlineInputBorder(),
                          ),
                          items:
                              Gender.values.map((gender) {
                                return DropdownMenuItem<Gender>(
                                  value: gender,
                                  child: Text(gender.toString()),
                                );
                              }).toList(),
                          onChanged: (newVal) {
                            setState(() {
                              chosenGender = newVal;
                            });
                          },
                        ),
                        SizedBox(height: 7),
                        DropdownButtonFormField<BowClass>(
                          value: chosenBow,
                          decoration: const InputDecoration(
                            labelText: 'Класс лука',
                            border: OutlineInputBorder(),
                          ),
                          items:
                              BowClass.values.map((bow) {
                                return DropdownMenuItem<BowClass>(
                                  value: bow,
                                  child: Text(bow.toString()),
                                );
                              }).toList(),
                          onChanged: (newVal) {
                            setState(() {
                              chosenBow = newVal;
                            });
                          },
                        ),
                        SizedBox(height: 7),
                        DropdownButtonFormField<SportsRank>(
                          value: chosenRank,
                          decoration: const InputDecoration(
                            labelText: 'Спортивный разряд',
                            isDense: true,
                            contentPadding: EdgeInsets.symmetric(
                              vertical: 12,
                              horizontal: 8,
                            ),
                            border: OutlineInputBorder(),
                          ),
                          items:
                              SportsRank.values.map((rank) {
                                return DropdownMenuItem<SportsRank>(
                                  value: rank,
                                  child: Text(
                                    rank.toString(),
                                    overflow: TextOverflow.ellipsis,
                                  ),
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
                        buildFederation(),
                        SizedBox(height: 7),
                        buildClub(),
                        SizedBox(height: 7),
                        SizedBox(
                          width: double.infinity,
                          child: OutlinedButton(
                            onPressed: () => chooseDate(context),
                            child: Text(
                              "Дата рождения: $birthDate",
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
                                shape: StadiumBorder(),
                              ),
                              child: const Text(
                                "Отменить",
                                style: TextStyle(color: Colors.white),
                              ),
                            ),
                            ElevatedButton(
                              onPressed: () async {
                                final competitor = ChangeCompetitor(
                                  fullNameController.text,
                                  birthDate,
                                  chosenGender!,
                                  chosenBow,
                                  chosenRank,
                                  regionController.text,
                                  federationController.text,
                                  clubController.text,
                                );

                                await userProvider.setUser(api, competitor);

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
                                shape: const StadiumBorder(),
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
    );
  }

  Future chooseDate(BuildContext context) async {
    DateTime initialDate;
    initialDate = DateTime.parse(birthDate!);

    final newDate = await showDatePicker(
      context: context,
      initialDate: initialDate,
      firstDate: DateTime(1950),
      lastDate: DateTime.now(),
    );

    if (newDate == null) return;

    setState(() {
      birthDate = formatDate(newDate);
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

  Widget buildFederation() {
    return TextFormField(
      controller: federationController,
      decoration: InputDecoration(
        label: Text("Федерация", style: TextStyle()),
        prefixIcon: Icon(Icons.people_alt),
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
