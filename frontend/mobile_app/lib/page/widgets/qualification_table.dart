import 'package:flutter/material.dart';
import 'package:mobile_app/api/responses.dart';

class MyQualificationTable extends StatelessWidget {
  final Section? section;

  const MyQualificationTable({super.key, required this.section});

  @override
  Widget build(BuildContext context) {
    if (section == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: Theme.of(context).dividerColor),
        borderRadius: BorderRadius.circular(8),
      ),
      child: SingleChildScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        scrollDirection: Axis.horizontal,
        child: DataTable(
          columns: const [
            DataColumn(
              label: Text(
                "Место",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              numeric: true,
            ),
            DataColumn(
              label: Text(
                "Участник",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
            ),
            DataColumn(
              label: Text(
                "Итог",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              numeric: true,
            ),
          ],
          rows: [
            DataRow(
              cells: [
                DataCell(Text(section!.place.toString())),
                DataCell(
                  Text(
                    section!.competitor.fullName,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                DataCell(Text(section!.total.toString())),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
