import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/page/widgets/User.dart';
import 'package:mobile_app/page/widgets/appbar_widget.dart';
import 'package:mobile_app/page/widgets/profile_widget.dart';
import 'package:mobile_app/page/widgets/text_box.dart';

void main() {
  runApp(const EditProfilePage());
}

class AuthPage extends StatelessWidget {
  const AuthPage({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        primarySwatch: Colors.green,
      ),
      home: LoginPage(),
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  _LoginPage createState() => _LoginPage();
}

class _LoginPage extends State<LoginPage> {
  final _sizeBlackText = const TextStyle(fontSize: 12, color: Colors.black);
  final _sizeWhiteText = const TextStyle(fontSize: 15, color: Colors.white);

  bool isPasswordVisible = true;

  final myController = TextEditingController();
  var text = '';


  @override
  void dispose() {
    // Clean up the controller when the widget is disposed.
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
                    onPressed: null,
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
  _ProfilePageState createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final user = UserPreferences.myUser;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData(primaryColor: Colors.green.shade300),
      home: Scaffold(
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
                    onPressed: () {},
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
    );
  }
}

class EditProfilePage extends StatefulWidget {
  const EditProfilePage({super.key});

  @override
  State<StatefulWidget> createState() => _EditProfilePage();
}

class _EditProfilePage extends State<EditProfilePage> {
  final loginController = TextEditingController();
  final numberController = TextEditingController();

  @override
  void initState() {
    super.initState();
    loginController.addListener(() => setState(() {}));
  }

  Widget buildNumber() {
    return TextFormField(
          controller: numberController,
          maxLength: 12,
          validator: (val) {
            if(numberController.text.length < 11) {
              return "Введите корректный номер телефона";
            }
            return null;
          },
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
    return TextFormField(
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
    return MaterialApp(
      theme: ThemeData(brightness: Brightness.light),
      //darkTheme: ThemeData(brightness: Brightness.dark),
      //themeMode: ThemeMode.system,
      debugShowCheckedModeBanner: false,
      home: KeyboardDismisser(
        gestures: [
          GestureType.onTap,
        ],
        child: GestureDetector(
          child: Scaffold(
            appBar: AppBar(
              title: const Text("Редактирование профиля", style: TextStyle(fontSize: 13, fontWeight: FontWeight.bold),),
              centerTitle: true,


            ),
            body: SingleChildScrollView(
              child: Container(
                padding: const EdgeInsets.all(8),
                child: Column(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Form(child: Column(
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
                        ],
                      )
                      ),
                    )

                  ],
                ),
              ),
            ),
          ),
        ),
      )
    );
  }

}
