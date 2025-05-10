import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:mobile_app/page/widgets/CompetitionField.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

class MainCompetitionPage extends StatefulWidget {
  const MainCompetitionPage({super.key});

  @override
  State<MainCompetitionPage> createState() => _MainCompetitionPage();
}

class _MainCompetitionPage extends State<MainCompetitionPage> {
  List<Competition>? _competitions;
  List<Cup>? _cups;
   int _idCup = 0;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  @override
  Widget build(BuildContext context) {
    var competitions = _competitions;
    return RefreshIndicator(
      onRefresh: _loadCompetitions,
      child: Scaffold(
        appBar: OnionBar("Соревнования", context),
        body: Center(
          child: ListView.builder(itemCount: competitions?.length, itemBuilder:(context, index) {
            final competition = competitions?[index];
            return CompetitionField(nameOfComp: competition!.stage.name, date: _formatCompetitionDate(competition.startDate, competition.endDate));
          })
        ),
      ),
    );
  }

  Future<void> _loadData() async {
    final api = context.read<Api>();
    try {
      final cups = await api.getCups();
      if (cups.isNotEmpty) {
        setState(() {
          _cups = cups;
          _idCup = cups.first.id;
        });
      }
      _loadCompetitions();
    } catch (e) {
      throw "$e";
    }
  }

  Future<void> _loadCompetitions() async {
    try {
      final competitions = await context.read<Api>().getCupsCompetitions(_idCup);
      setState(() {
        _competitions = competitions;
      });
    } catch (e) {
      throw "$e";
    }
  }

  String _formatCompetitionDate(DateTime? start, DateTime? end) {
    if (start == null && end == null) {
      return "соревнований не указаны";
    }
    if (start == null) return 'до ${_formatDate(end!)}';
    if (end == null) return 'с ${_formatDate(start)}';

    if (start.day == end.day && start.month == end.month &&
        start.year == end.year) {
      return _formatDate(start);
    }

    if (start.year == end.year && start.month == end.month) {
      return '${start.day}-${end.day} ${_getMonth(start.month)} ${start
          .year} г.';
    }

    if (start.year == end.year) {
      return '${start.day} ${_getMonth(start.month)} - '
          '${end.day} ${_getMonth(end.month)} ${start.year} г.';
    }

    return '${_formatDate(start)} - ${_formatDate(end)}';
  }

  String _getMonth(int month) {
    const months = ['января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
      'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря'];
    return months[month - 1];
  }

  String _formatDate(DateTime date) {
    return '${date.day} ${_getMonth(date.month)} ${date.year} г.';
  }
}
