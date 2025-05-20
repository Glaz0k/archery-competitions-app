import 'package:flutter/material.dart';
import 'package:mobile_app/model/user_model.dart';
import 'package:mobile_app/page/widgets/qualification_table.dart';
import 'package:provider/provider.dart';

import '../api/api.dart';
import '../api/responses.dart';

class QualificationPage extends StatefulWidget {
  final int groupId;

  const QualificationPage({super.key, required this.groupId});

  @override
  State<QualificationPage> createState() => _QualificationPageState();
}

class _QualificationPageState extends State<QualificationPage> {
  Section? _section;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    var api = context.read<Api>();
    int userId = context.read<UserModel>().getId()!;
    QualificationTable table = await api.getIndividualGroupQualificationTable(
      widget.groupId,
    );
    setState(() {
      _section = table.sections.firstWhere(
        (section) => section.competitor.id == userId,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: _loadData,
      child: SingleChildScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                "Секция ${_section?.id}",
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const SizedBox(height: 16),
              MyQualificationTable(section: _section),
            ],
          ),
        ),
      ),
    );
  }
}
