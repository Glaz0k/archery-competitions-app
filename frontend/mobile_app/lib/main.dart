
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/page/widgets/User.dart';
import 'package:mobile_app/page/widgets/appbar_widget.dart';
import 'package:mobile_app/page/widgets/profile_widget.dart';
import 'package:mobile_app/page/widgets/text_box.dart';
import 'package:mobile_app/page/widgets/MainCompetitionPage.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';

void main() {

  runApp(const AuthPage());
}

class AuthPage extends StatelessWidget {
  const AuthPage({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => UserProvider(),
      child: MaterialApp(
        debugShowCheckedModeBanner: false,
        initialRoute: '/',
        routes: {
          '/': (context) => const LoginPage(),
          '/profile_page': (context) => ProfilePage(),
          '/edit_profile_page': (context) => EditProfilePage(),
          '/competition_page': (context) => MainCompetitionPage(),
        },
        theme: ThemeData(primarySwatch: Colors.green),
      ),
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPage();
}

class _LoginPage extends State<LoginPage> {
  final _sizeBlackText = const TextStyle(fontSize: 12, color: Colors.black);
  final _sizeWhiteText = const TextStyle(fontSize: 15, color: Colors.white);

  bool isPasswordVisible = true;

  final myController = TextEditingController();
  var text = '';

  @override
  void dispose() {
    myController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: AppBar(
            title: Text('Вход'),
            backgroundColor: Colors.green,
            centerTitle: true,
          ),
          body: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                SizedBox(
                  width: 300,
                  child: TextFormField(
                    controller: myController,
                    decoration: InputDecoration(
                      label: Text(
                        "Логин",
                        style: TextStyle(color: Colors.green, fontSize: 20),
                      ),
                      prefixIcon: const Icon(Icons.person),
                      hintText: 'Введите логин',
                      enabledBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                    ),
                    validator: (String? value) {},
                    onChanged: null,
                  ),
                ),
                Container(
                  width: 300,
                  padding: EdgeInsets.only(top: 10.0),
                  child: TextFormField(
                    decoration: InputDecoration(
                      label: Text(
                        "Пароль",
                        style: TextStyle(color: Colors.green, fontSize: 20),
                      ),
                      hintText: 'Введите пароль',
                      prefixIcon: const Icon(Icons.person),
                      enabledBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      suffixIcon: IconButton(
                        onPressed:
                            () => setState(
                              () => isPasswordVisible = !isPasswordVisible,
                            ),
                        icon:
                            isPasswordVisible
                                ? Icon(Icons.visibility_off)
                                : Icon(Icons.visibility),
                      ),
                    ),

                    obscureText: isPasswordVisible,
                    style: _sizeBlackText,
                  ),
                ),
                Padding(
                  padding: EdgeInsets.only(top: 25.0),
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pushNamed(context, '/profile_page');
                    },
                    style: ButtonStyle(
                      minimumSize: WidgetStateProperty.all(const Size(200, 40)),
                      backgroundColor: WidgetStateProperty.all(Colors.green),
                      overlayColor: WidgetStateProperty.all(Colors.blueAccent),
                    ),
                    child: Text("Войти", style: _sizeWhiteText),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {

  String formatDate(DateTime date) {
    try {
      final format = DateFormat('dd.MM.yyyy', 'ru_RU');
      return format.format(date);
    } catch (e) {
      return '${date.day.toString().padLeft(2, '0')}.${date.month.toString().padLeft(2, '0')}.${date.year}';;
    }
  }

  @override
  Widget build(BuildContext context) {
    final user = Provider.of<UserProvider>(context).userPref;
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: buildAppBar(context),
          body: SingleChildScrollView(
            physics: AlwaysScrollableScrollPhysics(), //BouncingScrollPhysics
            child: Container(
              padding: const EdgeInsets.all(8),
              child: Column(
                children: [
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Фамилия Имя"),
                      subtitle: Text(user.fullName),
                      leading: Icon(
                        CupertinoIcons.person_alt,
                        color: Colors.teal,
                      ),
                    ),
                  ),

                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Дата рождения"),
                      subtitle: Text(formatDate(user.dateOfBirth)),
                      leading: Icon(Icons.date_range, color: Colors.teal),
                    ),
                  ),

                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Пол"),
                      subtitle: Text(user.identity.getGender),
                      leading: Icon(Icons.perm_identity, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Лук"),
                      subtitle: Text(user.bow!.getBowClass),
                      leading: Icon(Icons.dangerous, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Спортивный разряд"),
                      subtitle: Text(user.rank!.getSportsRank),
                      leading: Icon(Icons.star, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Город"),
                      subtitle: Text(user.region),
                      leading: Icon(Icons.location_city, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Клуб"),
                      subtitle: Text(user.club),
                      leading: Icon(
                        CupertinoIcons.sportscourt_fill,
                        color: Colors.teal,
                      ),
                    ),
                  ),
                  ProfileWidget(onClicked: () async {}),
                  SizedBox(height: 30),
                  SizedBox(
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.pushNamed(
                          context,
                          '/edit_profile_page',
                          arguments: user,
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.green,
                        side: BorderSide.none,
                        shape: StadiumBorder(),
                      ),
                      child: const Text(
                        "Редактировать",
                        style: TextStyle(color: Colors.white),
                      ),
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
}

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
          appBar: AppBar(
            title: const Text(
              "Редактирование профиля",
              style: TextStyle(fontSize: 13, fontWeight: FontWeight.bold),
            ),
            centerTitle: true,
          ),
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
                              items: ['Мужчина', 'Женщина'].map((String value) {
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
                              items: ['Классический', 'Блочный',"КЛ(новички)","3Д-классический лук","3Д-составной лук","3Д-длинный лук","Периферийный лук","Периферийный лук(с кольцом)"].map((String value) {
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
                              items: ['Мастер спорта', 'Кандидат мастера спорта',"Первый разряд","Второй разряд","Третий разряд","Международный магистр","Периферийный лук","Заслуженный мастер спорта"].map((String value) {
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
                              width: double.infinity, // Растягивает на всю ширину
                              child: OutlinedButton(
                                onPressed: () => chooseDate(context),
                                child: Text('Изменить дату рождения',style: TextStyle(fontSize: 13, fontWeight: FontWeight.bold),),
                              ),
                            ),
                            SizedBox(height: 7),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                              children: [
                                ElevatedButton(
                                  onPressed: () {
                                    Navigator.pushNamed(
                                      context,
                                      "/profile_page",
                                    );
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
                                    userProvider.updateRegion(regionController.text);
                                    userProvider.updateClub(clubController.text);
                                    chosenGender == "Мужчина" ? userProvider.updateGender(Gender.male) : userProvider.updateGender(Gender.female);
                                    userProvider.updateBow(BowExtension.setBowClass(chosenBow));
                                    userProvider.updateRank(SportsRankExtension.setSportsRank(chosenRank));
                                    userProvider.updateDateOfBirth(chosenTime);
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        behavior: SnackBarBehavior.floating,
                                        backgroundColor: Colors.green,
                                        content: Container(
                                          decoration: BoxDecoration(
                                            color: Colors.green,
                                          ),
                                          child: Row(
                                            children: [
                                              Icon(
                                                Icons.check_circle,
                                                color: Colors.white,
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
                                    Navigator.pushReplacementNamed(context, '/profile_page');
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

  // Widget buildLogin() {
  //   return TextFormField(
  //     controller: loginController,
  //     decoration: InputDecoration(
  //       label: Text("Логин", style: TextStyle()),
  //       prefixIcon: Icon(Icons.login_sharp),
  //       suffixIcon:
  //           loginController.text.isEmpty
  //               ? Container(width: 1)
  //               : IconButton(
  //                 onPressed: () => loginController.clear(),
  //                 icon: Icon(Icons.clear),
  //               ),
  //       border: OutlineInputBorder(),
  //     ),
  //     textInputAction: TextInputAction.done,
  //   );
  // }

  Future chooseDate(BuildContext context) async{
    final newDate = await showDatePicker(context: context,
      initialDate: chosenTime,
      firstDate: DateTime(1950),
      lastDate: DateTime.now(),);
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
