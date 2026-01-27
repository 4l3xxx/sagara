import { Controller, Get, Param } from '@nestjs/common';
import { ApiTags, ApiOperation } from '@nestjs/swagger';
import { WorksService } from './works.service';

@ApiTags('Works')
@Controller('works')
export class WorksController {
  constructor(private readonly service: WorksService) {}

  @Get()
  @ApiOperation({ summary: 'Get all works' })
  findAll() {
    return this.service.findAll();
  }

  @Get(':slug')
  @ApiOperation({ summary: 'Get work by slug' })
  findOne(@Param('slug') slug: string) {
    return this.service.findBySlug(slug);
  }
}
