import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/page/widgets/User.dart';
import 'package:mobile_app/page/widgets/appbar_widget.dart';
import 'package:mobile_app/page/widgets/profile_widget.dart';
import 'package:mobile_app/page/widgets/text_box.dart';

void main() {
  runApp(const AuthPage());
}

class AuthPage extends StatelessWidget {
  const AuthPage({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      initialRoute: '/',
      routes: {
        '/': (context) => const LoginPage(),
        '/profile_page': (context) => ProfilePage(),
        '/edit_profile_page': (context) => EditProfilePage(),
      },
      theme: ThemeData(
        primarySwatch: Colors.green,
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
      gestures: [
        GestureType.onTap,
      ],
      child: GestureDetector(
        child: Scaffold(
          appBar: AppBar(title: Text('Вход'), backgroundColor: Colors.green, centerTitle: true,),
          body: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                SizedBox(
                  width: 300,
                  child: TextFormField(
                    controller: myController,
                    decoration: InputDecoration(
                      label: Text("Логин", style: TextStyle(color: Colors.green, fontSize: 20),),
                      prefixIcon: const Icon(Icons.person),
                      hintText: 'Введите логин',
                      enabledBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      )
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
                        label: Text("Пароль", style: TextStyle(color: Colors.green, fontSize: 20)),
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
                        suffixIcon: IconButton(onPressed: () => setState(() => isPasswordVisible = !isPasswordVisible), icon: isPasswordVisible ? Icon(Icons.visibility_off) : Icon(Icons.visibility))
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
  final user = UserPreferences.myUser;

  @override
  Widget build(BuildContext context) {
    return KeyboardDismisser(
      gestures: [
        GestureType.onTap
      ],
      child: GestureDetector(
        child: Scaffold(
          appBar: buildAppBar(context),
          body: SingleChildScrollView(
            physics: BouncingScrollPhysics(),
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
                      title: Text("Фамилия"),
                      subtitle: Text(user.surname),
                      leading: Icon(CupertinoIcons.person_alt, color: Colors.teal),
                    ),
                  ),

                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: Text("Имя"),
                      subtitle: Text(user.name),
                      leading: Icon(Icons.account_circle_sharp, color: Colors.teal),
                    ),
                  ),

                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: Text("Отчество"),
                      subtitle: Text(user.middleName),
                      leading: Icon(Icons.account_circle_sharp, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: Text("Логин"),
                      subtitle: Text(user.login),
                      leading: Icon(Icons.login_sharp, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: Text("Номер телефона"),
                      subtitle: Text(user.phoneNumber),
                      leading: Icon(Icons.phone, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: Text("Город"),
                      subtitle: Text(user.city),
                      leading: Icon(Icons.location_city, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: Text("Клуб"),
                      subtitle: Text(user.club),
                      leading: Icon(CupertinoIcons.sportscourt_fill, color: Colors.teal),
                    ),
                  ),
                  ProfileWidget(onClicked: () async {}),
                  SizedBox(height: 30),
                  SizedBox(
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.pushNamed(context, '/edit_profile_page', arguments: user);
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
  var loginController = TextEditingController();
  var nameController = TextEditingController();
  var surnameController = TextEditingController();
  var middleNameController = TextEditingController();
  var numberController = TextEditingController();
  var user = UserPreferences.myUser;

  @override
  void initState() {
    super.initState();
    loginController.text = user.login;
    numberController.text = user.phoneNumber;
    nameController.text = user.name;
    surnameController.text = user.surname;
    middleNameController.text = user.middleName;
  }

  @override
  void dispose() {
    loginController.dispose();
    numberController.dispose();
    nameController.dispose();
    surnameController.dispose();
    middleNameController.dispose();
    super.dispose();
  }

  Widget buildNumber() {
    return TextFormField(
          controller: numberController,
          maxLength: 12,
          // validator: (val) {
          //   if(numberController.text.length < 11) {
          //     return "Введите корректный номер телефона";
          //   }
          //   return null;
          // },
          decoration: InputDecoration(
              label: Text("Номер телефона", style: TextStyle(color: Colors.green),),
              hintText: "Введите номер телефона",
              counterText: '',
              prefixIcon: Icon(Icons.phone_iphone, color: Colors.blueAccent,),
              enabledBorder: OutlineInputBorder(
                borderSide: BorderSide(color: Colors.teal, width: 2.2),
                borderRadius: BorderRadius.circular(15),
              ),
              focusedBorder: OutlineInputBorder(
                borderSide: BorderSide(color: Colors.teal, width: 2.2),
                borderRadius: BorderRadius.circular(15),
              )
          ),
          keyboardType: TextInputType.phone,
    );
  }

  @override
  Widget build(BuildContext context) {
    return buildEditProfilePage(context);
  }
  Widget buildLogin() {
    return TextFormField(
      controller: loginController,
      decoration: InputDecoration(

          label: Text("Логин", style: TextStyle(),),
          prefixIcon: Icon(Icons.login_sharp),
          suffixIcon: loginController.text.isEmpty ? Container(width: 1,) : IconButton(onPressed: () => loginController.clear(), icon: Icon(Icons.clear)),
          border: OutlineInputBorder(
          )
      ),
      textInputAction: TextInputAction.done,

    );
  }

  Widget buildFIO(String str) {
    TextEditingController controller;

    if (str == "Фамилия") {
      controller = surnameController;
    } else if (str == "Имя") {
      controller = nameController;
    } else {
      controller = middleNameController;
    }
    return TextFormField(
      controller: controller,
      decoration: InputDecoration(
          label: Text(str, style: TextStyle(),),
          prefixIcon: Icon(Icons.person),
          border: OutlineInputBorder(
          )
      ),
      textInputAction: TextInputAction.done,

    );
  }

  Widget buildEditProfilePage(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Редактирование профиля", style: TextStyle(fontSize: 13, fontWeight: FontWeight.bold),),
        centerTitle: true,
      ),
      body: KeyboardDismisser(
        gestures: [
          GestureType.onTap
        ],
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
                        buildFIO("Фамилия"),
                        SizedBox(height: 7,),
                        buildFIO("Имя"),
                        SizedBox(height: 7,),
                        buildFIO("Отчество"),
                        SizedBox(height: 7,),
                        buildLogin(),
                        SizedBox(height: 7,),
                        buildNumber(),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            ElevatedButton(onPressed: () {
                              Navigator.pushNamed(context, "/profile_page");
                            },
                              style: ElevatedButton.styleFrom(
                                backgroundColor: Colors.redAccent,
                                side: BorderSide.none,
                                shape: StadiumBorder(),
                              ), child: const Text("Отменить", style: TextStyle(color: Colors.white)),
                            ),
                            ElevatedButton(onPressed: () {
                            },
                              style: ElevatedButton.styleFrom(
                                backgroundColor: Colors.green,
                                side: BorderSide.none,
                                shape: StadiumBorder(),
                              ), child: const Text("Сохранить", style: TextStyle(color: Colors.white)),
                            ),
                          ],
                        )
                      ],
                    ),
                  )
                ],
              ),
            ),
          ),
        ),
      )
    );
  }
}
